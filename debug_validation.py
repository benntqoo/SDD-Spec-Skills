import sys
import os
sys.path.insert(0, os.path.join(os.getcwd(), 'skills-base', 'vibe-integrity'))

# Import the module
import validate_vibe_integrity

# Create an instance
validator = validate_vibe_integrity.VibeIntegrityValidator()

# Run the validation
success, message = validator.run_validation()

# Print the results
print("Success:", success)
print("Message:", message)
print("Errors:", validator.errors)
print("Warnings:", validator.warnings)
print("Passed:", validator.passed)
