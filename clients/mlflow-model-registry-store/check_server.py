#!/usr/bin/env python3
"""Check if Model Registry server is running and accessible."""

import requests
import sys


def check_server(host="127.0.0.1", port=8080):
    """Check if Model Registry server is accessible."""
    base_url = f"http://{host}:{port}"
    health_url = f"{base_url}/api/model_registry/v1alpha3/registered_models?pageSize=1"

    print(f"🔍 Checking Model Registry server at {base_url}")

    try:
        # Try to make a simple request to the API
        response = requests.get(health_url, timeout=5)

        if response.status_code == 200:
            print(f"✅ Model Registry server is running and accessible")
            print(f"   Status: {response.status_code}")
            print(f"   URL: {health_url}")
            return True
        else:
            print(f"⚠️  Model Registry server responded but with status: {response.status_code}")
            print(f"   Response: {response.text[:200]}...")
            return False

    except requests.exceptions.ConnectionError:
        print(f"❌ Cannot connect to Model Registry server at {base_url}")
        print(f"   Make sure the server is running and accessible")
        print(f"   Try: kubectl port-forward svc/model-registry-service 8080:8080")
        return False

    except requests.exceptions.Timeout:
        print(f"⏰ Timeout connecting to Model Registry server at {base_url}")
        return False

    except Exception as e:
        print(f"❌ Error checking server: {e}")
        return False


def main():
    """Main function."""
    if len(sys.argv) > 1:
        try:
            port = int(sys.argv[1])
        except ValueError:
            print("Usage: python check_server.py [port]")
            sys.exit(1)
    else:
        port = 8080

    success = check_server(port=port)

    if success:
        print("\n🎉 Ready to run integration tests!")
        print("   Run: uv run pytest tests/test_integration.py -v")
    else:
        print("\n💡 To start Model Registry server:")
        print("   kubectl port-forward svc/model-registry-service 8080:8080")
        print("   # Then try this script again")

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()