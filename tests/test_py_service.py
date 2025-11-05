"""Integration test to verify Python service builds correctly."""

import os
import sys


def test_python_service_builds():
    """Test that the Python service builds successfully."""
    # The binary is in the runfiles because it's declared as data dependency
    server_path = "py_service/server"

    if os.path.exists(server_path):
        print(f"✓ Python service binary exists at {server_path}")
        return 0
    else:
        print("✗ Python service binary not found")
        print(f"Looking for: {server_path}")
        print("Current directory:", os.getcwd())
        print("Available files:")
        for root, dirs, files in os.walk("."):
            for file in files:
                if "server" in file:
                    print(f"  {os.path.join(root, file)}")
        return 1


if __name__ == "__main__":
    sys.exit(test_python_service_builds())
