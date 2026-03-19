#!/bin/bash
# ============================================================
# VIC-SDD Periodic Verification Script
# ============================================================
# 
# 这个脚本用于外部验证 AI 的声明
# 可以由人类手动运行，或定时触发
#
# 用法：
#   ./scripts/verify.sh
#   ./scripts/verify.sh --quiet  # 简化输出
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 参数解析
QUIET=false
if [ "$1" == "--quiet" ]; then
    QUIET=true
fi

# 项目路径
PROJECT_DIR=$(pwd)
VIC_DIR="${PROJECT_DIR}/.vic-sdd"

# ============================================================
# 辅助函数
# ============================================================

log() {
    if [ "$QUIET" == "false" ]; then
        echo -e "$1"
    fi
}

check_pass() {
    log "   ${GREEN}✓${NC} $1"
}

check_fail() {
    log "   ${RED}✗${NC} $1"
}

check_warn() {
    log "   ${YELLOW}⚠${NC} $1"
}

check_info() {
    log "   ${BLUE}ℹ${NC} $1"
}

# ============================================================
# 验证开始
# ============================================================

log ""
log "═══════════════════════════════════════════════════════════"
log "  🔍 VIC-SDD Verification"
log "═══════════════════════════════════════════════════════════"
log ""

# 检查 VIC 目录是否存在
if [ ! -d "$VIC_DIR" ]; then
    log "${RED}VIC directory not found: $VIC_DIR${NC}"
    exit 1
fi

# ============================================================
# 1. 知识边界检查
# ============================================================

log "📋 Knowledge Boundary Check"
log "───────────────────────────────────────────────────────────"

KB_FILE="${VIC_DIR}/knowledge-boundary.yaml"
if [ -f "$KB_FILE" ]; then
    # 检查未验证的推断
    unverified=$(grep -c "needs_verification: true" "$KB_FILE" 2>/dev/null || echo "0")
    if [ "$unverified" -gt 5 ]; then
        check_warn "Unverified inferences: ${unverified} (建议先验证)"
    else
        check_pass "Unverified inferences: ${unverified}"
    fi
    
    # 检查高风险假设
    high_risk=$(grep -c "risk: high" "$KB_FILE" 2>/dev/null || echo "0")
    if [ "$high_risk" -gt 2 ]; then
        check_warn "High-risk assumptions: ${high_risk} (需要验证)"
    else
        check_pass "High-risk assumptions: ${high_risk}"
    fi
    
    # 检查阻塞未知
    blocking=$(grep -c "blocks:" "$KB_FILE" 2>/dev/null || echo "0")
    if [ "$blocking" -gt 0 ]; then
        check_info "Blocking unknowns: ${blocking}"
    fi
else
    check_warn "knowledge-boundary.yaml not found"
fi

# ============================================================
# 2. 信号检查
# ============================================================

log ""
log "📋 Signal Register Check"
log "───────────────────────────────────────────────────────────"

SR_FILE="${VIC_DIR}/signal-register.yaml"
if [ -f "$SR_FILE" ]; then
    # 检查阻塞信号
    blockers=$(grep -c "severity: high" "$SR_FILE" 2>/dev/null || echo "0")
    blockers="${blockers//[^0-9]/}"
    [ -z "$blockers" ] && blockers=0
    if [ "$blockers" -gt 2 ]; then
        check_fail "High-severity blockers: ${blockers} (超过阈值)"
    elif [ "$blockers" -gt 0 ]; then
        check_warn "High-severity blockers: ${blockers}"
    else
        check_pass "High-severity blockers: 0"
    fi
    
    # 检查信号新鲜度
    last_signal=$(grep "timestamp:" "$SR_FILE" 2>/dev/null | tail -1 || echo "")
    if [ -n "$last_signal" ]; then
        check_info "Last signal: ${last_signal#timestamp: }"
    fi
    
    # 检查信心度
    confidence=$(grep "calculated:" "$SR_FILE" 2>/dev/null | tail -1 | grep -oE "[0-9]+\.[0-9]+" || echo "0")
    confidence="${confidence//[^0-9.]/}"
    [ -z "$confidence" ] && confidence=0
    
    # 简单比较（不依赖 bc）
    if [ "$(echo "$confidence < 0.4" | awk '{print ($1 < $3) ? "1" : "0"}')" = "1" ]; then
        check_fail "Confidence: ${confidence} (过低，建议暂停)"
    elif [ "$(echo "$confidence < 0.7" | awk '{print ($1 < $3) ? "1" : "0"}')" = "1" ]; then
        check_warn "Confidence: ${confidence} (中等)"
    else
        check_pass "Confidence: ${confidence}"
    fi
else
    check_warn "signal-register.yaml not found"
fi

# ============================================================
# 3. 质量红线检查
# ============================================================

log ""
log "📋 Quality Check"
log "───────────────────────────────────────────────────────────"

# 检查 TODO/FIXME
todos=$(grep -r "TODO\|FIXME" "${PROJECT_DIR}/src" 2>/dev/null | grep -v ".test." | grep -v "node_modules" | wc -l)
if [ "$todos" -gt 0 ]; then
    check_fail "TODO/FIXME in src/: ${todos}"
    grep -r "TODO\|FIXME" "${PROJECT_DIR}/src" 2>/dev/null | grep -v ".test." | grep -v "node_modules" | head -3 | while read line; do
        check_info "  $line"
    done
else
    check_pass "TODO/FIXME in src/: 0"
fi

# 检查 console.log
consoles=$(grep -r "console.log" "${PROJECT_DIR}/src" 2>/dev/null | grep -v ".test." | grep -v "node_modules" | grep -v "//.*console.log" | wc -l)
if [ "$consoles" -gt 0 ]; then
    check_fail "console.log in src/: ${consoles}"
else
    check_pass "console.log in src/: 0"
fi

# 检查硬编码密钥
secrets=$(grep -r "password\s*=\s*[\"'][^\"']*[\"']" "${PROJECT_DIR}/src" 2>/dev/null | grep -v ".test." | grep -v "node_modules" | grep -v "placeholder" | wc -l)
if [ "$secrets" -gt 0 ]; then
    check_fail "Potential hardcoded secrets found"
else
    check_pass "No hardcoded secrets detected"
fi

# ============================================================
# 4. SPEC 对齐检查
# ============================================================

log ""
log "📋 SPEC Alignment Check"
log "───────────────────────────────────────────────────────────"

SPEC_FILE="${VIC_DIR}/SPEC-REQUIREMENTS.md"
if [ -f "$SPEC_FILE" ]; then
    # 统计验收标准
    acceptance=$(grep -c "^\s*-\s*\[" "$SPEC_FILE" 2>/dev/null || echo "0")
    check_info "Acceptance criteria: ${acceptance}"
    
    # 统计测试（如果有 tests 目录）
    if [ -d "${PROJECT_DIR}/tests" ] || [ -d "${PROJECT_DIR}/test" ]; then
        test_dir="${PROJECT_DIR}/tests"
        [ ! -d "$test_dir" ] && test_dir="${PROJECT_DIR}/test"
        tests=$(grep -c "describe\|it(" "${test_dir}" 2>/dev/null | head -1 || echo "0")
        check_info "Tests found: ${tests}"
        
        if [ "$tests" -lt "$acceptance" ]; then
            check_warn "Tests < Acceptance criteria (${tests}/${acceptance})"
        else
            check_pass "Tests >= Acceptance criteria (${tests}/${acceptance})"
        fi
    fi
else
    check_warn "SPEC-REQUIREMENTS.md not found"
fi

# ============================================================
# 5. 探索日志检查
# ============================================================

log ""
log "📋 Exploration Journal Check"
log "───────────────────────────────────────────────────────────"

EJ_FILE="${VIC_DIR}/exploration-journal.yaml"
if [ -f "$EJ_FILE" ]; then
    # 检查失败尝试
    failed=$(grep -c "result: failed" "$EJ_FILE" 2>/dev/null || echo "0")
    check_info "Failed attempts recorded: ${failed}"
    
    # 检查最新探索
    last_explore=$(grep "timestamp:" "$EJ_FILE" 2>/dev/null | tail -1 || echo "")
    if [ -n "$last_explore" ]; then
        check_info "Last exploration: ${last_explore#timestamp: }"
    fi
else
    check_warn "exploration-journal.yaml not found"
fi

# ============================================================
# 6. 决策记录检查
# ============================================================

log ""
log "📋 Decision Records Check"
log "───────────────────────────────────────────────────────────"

DG_FILE="${VIC_DIR}/decision-guardrails.yaml"
if [ -f "$DG_FILE" ]; then
    # 检查决策
    decisions=$(grep -c "decision:" "$DG_FILE" 2>/dev/null || echo "0")
    check_info "Decisions recorded: ${decisions}"
    
    # 检查停止的决策
    stopped=$(grep -c "result: stopped" "$DG_FILE" 2>/dev/null || echo "0")
    if [ "$stopped" -gt 0 ]; then
        check_warn "Stopped decisions: ${stopped}"
    fi
else
    check_warn "decision-guardrails.yaml not found"
fi

# ============================================================
# 总结
# ============================================================

log ""
log "═══════════════════════════════════════════════════════════"

# 检查是否有严重问题
if grep -q "✗ Verification Failed" <<< "$(cat)"; then
    log "  ${RED}✗ Verification Failed - 需要处理问题${NC}"
elif grep -q "⚠" <<< "$(cat)"; then
    log "  ${YELLOW}⚠ Verification Passed with Warnings${NC}"
else
    log "  ${GREEN}✓ Verification Passed${NC}"
fi

log "═══════════════════════════════════════════════════════════"
log ""

# 输出帮助信息
if [ "$FAIL_COUNT" -gt 0 ]; then
    log "建议操作："
    log "  1. 查看 .vic-sdd/signal-register.yaml 的 blockers"
    log "  2. 检查 knowledge-boundary.yaml 的高风险项"
    log "  3. 运行 knowledge-boundary skill 验证推断"
    log "  4. 使用 pre-decision-check skill 重新检查"
fi

log ""
