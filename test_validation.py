#!/usr/bin/env python3
import sys
import os
print("Python version:", sys.version)
print("Current directory:", os.getcwd())
print("Files in current dir:", os.listdir('.'))

# Try to import the validation script directly
try:
    sys.path.insert(0, 'skills-base/vibe-integrity')
    import validate_vibe_integrity
    print("Import successful")
    print("Available attributes:", [x for x in dir(validate_vibe_integrity) if not x.startswith('_')])
    
    # Try to create validator
    validator = validate_vibe_integrity.VibeIntegrityValidator()
    print("Validator created")
    
    # Run validation
    success, msg = validator.run_validation()
    print(f"Validation result: success={success}, message='{msg}'")
    
except Exception as e:
    print(f"Error: {e}")
    import traceback
    traceback.print_exc()
