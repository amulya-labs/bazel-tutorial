"""Integration test to verify Python service builds correctly."""

import subprocess
import sys


def test_python_service_builds():
    """Test that the Python service builds successfully."""
    try:
        result = subprocess.run(
            ["bazel", "build", "//py_service:server"],
            capture_output=True,
            text=True,
            check=True
        )
        print("✓ Python service builds successfully")
        return 0
    except subprocess.CalledProcessError as e:
        print("✗ Python service failed to build")
        print(e.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(test_python_service_builds())
