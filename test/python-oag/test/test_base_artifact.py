# coding: utf-8

"""
    Model Registry REST API

    REST API for Model Registry to create and manage ML model metadata

    The version of the OpenAPI document: v1alpha3
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


import unittest
import datetime

from mr_openapi.models.base_artifact import BaseArtifact  # noqa: E501

class TestBaseArtifact(unittest.TestCase):
    """BaseArtifact unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> BaseArtifact:
        """Test BaseArtifact
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # uncomment below to create an instance of `BaseArtifact`
        """
        model = BaseArtifact()  # noqa: E501
        if include_optional:
            return BaseArtifact(
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
            return BaseArtifact(
        )
        """

    def testBaseArtifact(self):
        """Test BaseArtifact"""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
