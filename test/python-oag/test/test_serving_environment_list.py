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

from mr_openapi.models.serving_environment_list import ServingEnvironmentList  # noqa: E501

class TestServingEnvironmentList(unittest.TestCase):
    """ServingEnvironmentList unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> ServingEnvironmentList:
        """Test ServingEnvironmentList
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # uncomment below to create an instance of `ServingEnvironmentList`
        """
        model = ServingEnvironmentList()  # noqa: E501
        if include_optional:
            return ServingEnvironmentList(
                next_page_token = '',
                page_size = 56,
                size = 56,
                items = [
                    null
                    ]
            )
        else:
            return ServingEnvironmentList(
                next_page_token = '',
                page_size = 56,
                size = 56,
        )
        """

    def testServingEnvironmentList(self):
        """Test ServingEnvironmentList"""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
