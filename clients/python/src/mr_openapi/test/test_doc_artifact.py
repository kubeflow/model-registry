"""Model Registry REST API.

REST API for Model Registry to create and manage ML model metadata

The version of the OpenAPI document: v1alpha3
Generated by OpenAPI Generator (https://openapi-generator.tech)

Do not edit the class manually.
"""  # noqa: E501

import unittest

from mr_openapi.models.doc_artifact import DocArtifact


class TestDocArtifact(unittest.TestCase):
    """DocArtifact unit test stubs."""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> DocArtifact:
        """Test DocArtifact
        include_option is a boolean, when False only required
        params are included, when True both required and
        optional params are included.
        """
        # uncomment below to create an instance of `DocArtifact`
        """
        model = DocArtifact()
        if include_optional:
            return DocArtifact(
                artifact_type = 'doc-artifact',
                custom_properties = {
                    'key' : null
                    },
                description = '',
                external_id = '',
                uri = '',
                state = 'UNKNOWN',
                name = '',
                id = '',
                create_time_since_epoch = '',
                last_update_time_since_epoch = ''
            )
        else:
            return DocArtifact(
                artifact_type = 'doc-artifact',
        )
        """

    def testDocArtifact(self):
        """Test DocArtifact."""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)


if __name__ == "__main__":
    unittest.main()
