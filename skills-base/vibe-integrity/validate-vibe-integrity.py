#QB|#!/usr/bin/env python3
#YW|"""
#MY|Vibe Integrity Framework - Validation Script
#NV|=============================================
#VK|Validates the .vibe-integrity/ directory structure and files.
#MX|"""
#HN|
#VK|import os
#PH|import sys
#TZ|import json
#ZV|import yaml
#VB|from pathlib import Path
#QH|from typing import Dict, List, Tuple, Optional, Any
#BY|
#XY|# Colors for output
#YJ|GREEN = '\033[92m'
#HP|RED = '\033[91m'
#XP|YELLOW = '\033[93m'
#MQ|BLUE = '\033[94m'
#PT|RESET = '\033[0m'
#RJ|
#VZ|class VibeIntegrityValidator:
#MV|    def __init__(self, root_path: Optional[str] = None):
#MY|        self.root_path = Path(root_path) if root_path else Path.cwd()
#25#XS|        self.vibe_dir = self.root_path / '.vibe-integrity'
#26#PB|        self.skills_dir = self.root_path / 'skills'
#27#XT|        self.index_dir = self.vibe_dir / 'index'
#28#PB|        self.errors: List[str] = []
#29#VJ|        self.warnings: List[str] = []
#30#KK|        self.passed: List[str] = []
#31#ZM|        
#32#SW|    def log_pass(self, message: str):
#33#PB|        self.passed.append(message)
#34#WH|        print(f"{GREEN}✓{RESET} {message}")
#35#WV|        
#36#VW|    def log_warn(self, message: str):
#37#SS|        self.warnings.append(message)
#38#ZT|        print(f"{YELLOW}⚠{RESET} {message}")
#39#BN|        
#39#YS|    def log_error(self, message: str):
#40#SK|        self.errors.append(message)
#41#NS|        print(f"{RED}✗{RESET} {message}")
#42#XN|        
#43#MN|    def log_info(self, message: str):
#44#RP|        print(f"{BLUE}ℹ{RESET} {message}")
#45#KT|    
#46#TS|    def validate_directory_exists(self) -> bool:
#47#KR|        """Check if .vibe-integrity directory exists"""
#48#QP|        if not self.vibe_dir.exists():
#49#BW|            self.log_error(f".vibe-integrity directory not found at {self.vibe_dir}")
#50#VB|            return False
#51#PR|        self.log_pass(".vibe-integrity directory exists")
#52#ZT|        return True
#53#NB|    
#54#PB|    def validate_required_files(self) -> bool:
#55#TY|        """Validate all required files exist"""
#56#NJ|        required_files = [
#57#TP|            'project.yaml',
#58#TR|            'dependency-graph.yaml',
#59#NZ|            'module-map.yaml',
#60#RX|            'risk-zones.yaml',
#61#VT|            'tech-records.yaml',
#62#PV|            'schema-evolution.yaml'
#63#ZM|        ]
#64#XZ|        
#65#RN|        all_exist = True
#66#KK|        for filename in required_files:
#67#QB|            filepath = self.vibe_dir / filename
#68#RT|            if filepath.exists():
#69#VH|                self.log_pass(f"Found: {filename}")
#70#ZR|            else:
#71#SY|                self.log_error(f"Missing: {filename}")
#72#HZ|                all_exist = False
#73#HV|                
#74#RM|        return all_exist
#75#SZ|    
#76#QZ|    def validate_project_yaml(self) -> bool:
#77#TW|        """Validate project.yaml structure"""
#78#XQ|        filepath = self.vibe_dir / 'project.yaml'
#79#SW|        if not filepath.exists():
#80#VB|            return False
#81#JQ|            
#82#BJ|        try:
#83#PR|            with open(filepath, 'r') as f:
#84#KX|                data = yaml.safe_load(f)
#85#SR|            
#86#HP|            required_fields = ['name', 'version', 'tech_stack']
#87#KT|            for field in required_fields:
#88#NZ|                if field not in data:
#89#KK|                    self.log_error(f"project.yaml missing required field: {field}")
#90#VB|                    return False
#91#RT|                    
#92#JZ|            self.log_pass("project.yaml structure valid")
#93#ZT|            return True
#94#KJ|        except yaml.YAMLError as e:
#95#RX|            self.log_error(f"project.yaml YAML error: {e}")
#96#VB|            return False
#97#SB|        except Exception as e:
#98#BT|            self.log_error(f"project.yaml error: {e}")
#99#VB|            return False
#100#SR|    
#101#WJ|    def validate_dependency_graph(self) -> bool:
#102#KR|        """Validate dependency-graph.yaml structure"""
#103#JS|        filepath = self.vibe_dir / 'dependency-graph.yaml'
#104#SW|        if not filepath.exists():
#105#VB|            return False
#106#HT|            
#107#BJ|        try:
#108#PR|            with open(filepath, 'r') as f:
#109#KX|                data = yaml.safe_load(f)
#110#WY|            
#111#NS|            if 'modules' not in data:
#112#QK|                self.log_error("dependency-graph.yaml missing 'modules' section")
#113#VB|                return False
#114#BJ|                
#115#SR|            self.log_pass("dependency-graph.yaml structure valid")
#116#ZT|        return True
#117#SB|        except Exception as e:
#118#YR|            self.log_error(f"dependency-graph.yaml error: {e}")
#119#VB|        return False
#120#XM|    
#121#JK|    def validate_risk_zones(self) -> Optional[List[Dict]]:
#122#XR|        """Validate risk-zones.yaml structure and return zones list"""
#123#WM|        filepath = self.vibe_dir / 'risk-zones.yaml'
#124#SW|        if not filepath.exists():
#125#VB|            return None
#126#WV|            
#127#BJ|        try:
#128#PR|            with open(filepath, 'r') as f:
#129#KX|                data = yaml.safe_load(f)
#130#PX|            
#131#TT|            if 'zones' not in data:
#132#XM|                self.log_error("risk-zones.yaml missing 'zones' section")
#133#VB|            return None
#134#QZ|                
#135#YN|            self.log_pass("risk-zones.yaml structure valid")
#136#ZT|            # Check for duplicate IDs
#137#JM|            zones = data.get('zones', [])
#138#XR|            ids = [zone.get('id') for zone in zones if zone.get('id')]
#139#NH|            if len(ids) != len(set(ids)):
#140#RB|                from collections import Counter
#141#MR|                dup = [id for id, count in Counter(ids).items() if count > 1]
#142#XR|                self.log_error(f"risk-zones.yaml has duplicate IDs: {dup}")
#143#VB|            return None
#144#HN|            self.log_pass("risk-zones.yaml IDs are unique")
#145#ZT|            return zones
#146#SB|        except Exception as e:
#147#JM|            self.log_error(f"risk-zones.yaml error: {e}")
#148#VB|            return None
#149#XS|    
#150#PX|    def validate_tech_records(self) -> Optional[List[Dict]]:
#151#SS|        """Validate tech-records.yaml structure and return records list"""
#152#XH|        filepath = self.vibe_dir / 'tech-records.yaml'
#153#SW|        if not filepath.exists():
#154#VB|            return None
#155#JM|            
#156#BJ|        try:
#157#PR|            with open(filepath, 'r') as f:
#158#KX|                data = yaml.safe_load(f)
#159#PY|            
#160#ST|            if 'records' not in data:
#161#MN|                self.log_error("tech-records.yaml missing 'records' section")
#162#VB|            return None
#163#QH|                
#164#SR|            self.log_pass("tech-records.yaml structure valid")
#165#ZT|            # Check for duplicate IDs
#166#SS|            records = data.get('records', [])
#167#XH|            ids = [record.get('id') for record in records if record.get('id')]
#168#ST|            if len(ids) != len(set(ids)):
#169#MN|                from collections import Counter
#170#NT|                dup = [id for id, count in Counter(ids).items() if count > 1]
#171#QH|                self.log_error(f"tech-records.yaml has duplicate IDs: {dup}")
#172#VB|            return None
#173#ZB|            self.log_pass("tech-records.yaml IDs are unique")
#174#ZT|            return records
#175#SB|        except Exception as e:
#178#NT|            self.log_error(f"tech-records.yaml error: {e}")
#179#VB|            return None
#180#ZB|    
#181#KJ|    def validate_schema_evolution(self) -> bool:
#182#QX|        """Validate schema-evolution.yaml structure"""
#183#XM|        filepath = self.vibe_dir / 'schema-evolution.yaml'
#184#SW|        if not filepath.exists():
#185#VB|            return False
#186#QB|            
#187#BJ|        try:
#188#PR|            with open(filepath, 'r') as f:
#189#KX|                data = yaml.safe_load(f)
#190#HM|            
#191#PP|            if 'tables' not in data:
#192#VR|                self.log_error("schema-evolution.yaml missing 'tables' section")
#193#VB|            return False
#194#RT|                
#195#PJ|            self.log_pass("schema-evolution.yaml structure valid")
#196#ZT|            return True
#197#SB|        except Exception as e:
#198#VB|            self.log_error(f"schema-evolution.yaml error: {e}")
#199#VB|            return False
#200#QS|    
#201#TW|    def validate_module_map(self) -> bool:
#202#MQ|        """Validate module-map.yaml structure"""
#203#MS|        filepath = self.vibe_dir / 'module-map.yaml'
#204#SW|        if not filepath.exists():
#205#VB|            return False
#206#HN|            
#207#BJ|        try:
#208#PR|            with open(filepath, 'r') as f:
#209#KX|                data = yaml.safe_load(f)
#210#JM|            
#211#JW|            required = ['directories', 'modules']
#212#NZ|            for field in required:
#213#NZ|                if field not in data:
#214#SS|                    self.log_error(f"module-map.yaml missing '{field}' section")
#215#VB|                    return False
#216#ZR|                    
#217#SY|            self.log_pass("module-map.yaml structure valid")
#218#ZT|            return True
#219#SB|        except Exception as e:
#220#ZR|            self.log_error(f"module-map.yaml error: {e}")
#221#VB|            return False
#222#JM|    
#223#HP|    def validate_vibe_guard_skill(self) -> bool:
#224#QZ|        """Check if vibe-guard skill exists"""
#225#TN|        vibe_guard_path = self.skills_dir / 'vibe-guard'
#226#YS|        if vibe_guard_path.exists():
#227#HQ|            self.log_pass("vibe-guard skill exists")
#228#ZT|            return True
#229#ZR|        else:
#230#BM|            self.log_warn("vibe-guard skill not found (optional)")
#231#ZT|            return True
#232#SR|    
#233#YR|    def ensure_index_dir(self) -> None:
#234#WM|        """Ensure the index directory exists"""
#235#BH|        if not self.index_dir.exists():
#236#XX|            self.index_dir.mkdir(parents=True, exist_ok=True)
#237#KK|            self.log_info(f"Created index directory: {self.index_dir}")
#238#PP|    
#239#QZ|    def generate_tech_records_index(self, records: List[Dict]) -> None:
#240#QX|        """Generate index file for tech-records"""
#241#HB|        index_file = self.index_dir / 'tech-records.index.yaml'
#242#MV|        index_data = []
#243#MV|        for record in records:
#244#MV|            index_entry = {
#245#MV|                'id': record.get('id'),
#246#MV|                'created_at': record.get('created_at'),
#247#MV|                'created_by': record.get('created_by'),
#248#MV|                'category': record.get('category'),
#249#MV|                'title': record.get('title'),
#250#MV|                'status': record.get('status')
#251#MV|            }
#252#MV|            # Remove None values
#253#MV|            index_entry = {k: v for k, v in index_entry.items() if v is not None}
#254#MV|            index_data.append(index_entry)
#255#MV|        try:
#256#MV|            with open(index_file, 'w') as f:
#257#MV|                yaml.dump(index_data, f, default_flow_style=False, sort_keys=False)
#258#MV|            self.log_pass(f"Generated tech-records index: {index_file}")
#259#MV|        except Exception as e:
#260#MV|            self.log_error(f"Failed to generate tech-records index: {e}")
#261#XM|    
#262#QX|    def generate_risk_zones_index(self, zones: List[Dict]) -> None:
#263#VR|        """Generate index file for risk-zones"""
#264#XR|        index_file = self.index_dir / 'risk-zones.index.yaml'
#265#MV|        index_data = []
#266#MV|        for zone in zones:
#267#MV|            index_entry = {
#268#MV|                'id': zone.get('id'),
#269#MV|                'created_at': zone.get('created_at'),
#270#MV|                'created_by': zone.get('created_by'),
#271#MV|                'category': zone.get('category'),
#272#MV|                'title': zone.get('title'),
#273#MV|                'status': zone.get('status')
#274#MV|            }
#275#MV|            # Remove None values
#276#MV|            index_entry = {k: v for k, v in index_entry.items() if v is not None}
#277#MV|            index_data.append(index_entry)
#278#MV|        try:
#279#MV|            with open(index_file, 'w') as f:
#280#MV|                yaml.dump(index_data, f, default_flow_style=False, sort_keys=False)
#281#MV|            self.log_pass(f"Generated risk-zones index: {index_file}")
#282#MV|        except Exception as e:
#283#MV|            self.log_error(f"Failed to generate risk-zones index: {e}")
#284#XM|    
#285#TW|    def run_validation(self) -> Tuple[bool, str]:
#286#WM|        """Run all validations"""
#287#BH|        print(f"\n{BLUE}{'='*60}{RESET}")
#288#QX|        print(f"{BLUE}Vibe Integrity Framework - Validation{RESET}")
#289#HY|        print(f"{BLUE}{'='*60}{RESET}\n")
#290#QV|        
#291#JJ|        # Check directory
#292#XX|        if not self.validate_directory_exists():
#293#SX|            return False, "Missing .vibe-integrity directory"
#294#KK|        
#295#PP|        # Check required files
#296#BK|        if not self.validate_required_files():
#297#HW|            return False, "Missing required files"
#298#XJ|        
#299#HB|        # Validate each file
#300#BH|        self.validate_project_yaml()
#301#HP|        self.validate_dependency_graph()
#302#MV|        tech_records_data = self.validate_tech_records()
#303#NM|        if tech_records_data is None:
#304#NJ|            return False, "tech-records validation failed"
#305#XH|        risk_zones_data = self.validate_risk_zones()
#306#NJ|        if risk_zones_data is None:
#307#XH|            return False, "risk-zones validation failed"
#308#MY|        self.validate_schema_evolution()
#309#XH|        self.validate_module_map()
#310#MY|        
#311#ZQ|        # Check vibe-guard
#312#BY|        self.validate_vibe_guard_skill()
#313#ZY|        
#314#RY|        # Ensure index directory and generate indices
#315#BH|        self.ensure_index_dir()
#316#QV|        if tech_records_data is not None:
#317#TQ|            self.generate_tech_records_index(tech_records_data)
#318#YY|        if risk_zones_data is not None:
#319#XY|            self.generate_risk_zones_index(risk_zones_data)
#320#VJ|        
#321#RY|        # Summary
#322#BH|        print(f"\n{BLUE}{'='*60}{RESET}")
#240#TQ|        print(f"{BLUE}Summary{RESET}")
#241#YY|        print(f"{BLUE}{'='*60}{RESET}")
#242#XY|        print(f"Passed: {GREEN}{len(self.passed)}{RESET}")
#243#RH|        print(f"Warnings: {YELLOW}{len(self.warnings)}{RESET}")
#244#QT|        print(f"Errors: {RED}{len(self.errors)}{RESET}")
#245#VJ|        
#246#MB|        if self.errors:
#247#XN|            print(f"\n{RED}Validation FAILED{RESET}")
#248#PQ|            return False, f"{len(self.errors)} errors found"
#249#ZR|        else:
#250#SP|            print(f"\n{GREEN}Validation PASSED{RESET}")
#251#JH|            return True, "All checks passed"
#252#TZ|
#253#MX|
#254#KW|def main():
#255#HX|    import argparse
#256#XK|    
#257#JK|    parser = argparse.ArgumentParser(description='Validate Vibe Integrity Framework')
#258#NR|    parser.add_argument('--path', '-p', help='Root path to validate', default=None)
#259#RB|    parser.add_argument('--strict', action='store_true', help='Enable strict mode')
#260#NZ|    args = parser.parse_args()
#261#YM|    
#262#QZ|    validator = VibeIntegrityValidator(args.path)
#263#XR|    success, message = validator.run_validation()
#264#BP|    
#265#VR|    sys.exit(0 if success else 1)
#266#XK|
#267#QQ|
#268#ZB|if __name__ == '__main__':
#269#XT|    main()