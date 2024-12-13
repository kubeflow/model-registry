module.exports = {
  "parser": "@typescript-eslint/parser",
  "env": {
    "browser": true,
    "node": true
  },
  // tell the TypeScript parser that we want to use JSX syntax
  "parserOptions": {
    "tsx": true,
    "jsx": true,
    "js": true,
    "useJSXTextNode": true,
    "project": "./tsconfig.json",
    "tsconfigRootDir": __dirname
  },
  // includes the typescript specific rules found here: https://github.com/typescript-eslint/typescript-eslint/tree/master/packages/eslint-plugin#supported-rules
  "plugins": [
    "@typescript-eslint",
    "react-hooks",
    "eslint-plugin-react-hooks",
    "import",
    "no-only-tests",
    "no-relative-import-paths",
    "prettier"
  ],
  "extends": [
    "eslint:recommended",
    "plugin:jsx-a11y/recommended",
    "plugin:react/recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:prettier/recommended",
    "prettier"
  ],
  "globals": {
    "window": "readonly",
    "describe": "readonly",
    "test": "readonly",
    "expect": "readonly",
    "it": "readonly",
    "process": "readonly",
    "document": "readonly"
  },
  "settings": {
    "react": {
      "version": "detect"
    },
    "import/parsers": {
      "@typescript-eslint/parser": [".ts", ".tsx"]
    },
    "import/resolver": {
      "typescript": {
        "project": "."
      }
    }
  },
  "rules": {
    "jsx-a11y/no-autofocus": ["error", { "ignoreNonDOM": true }],
    "jsx-a11y/anchor-is-valid": [
      "error",
      {
        "components": ["Link"],
        "specialLink": ["to"],
        "aspects": ["noHref", "invalidHref", "preferButton"]
      }
    ],
    "react/jsx-boolean-value": "error",
    "react/jsx-fragments": "error",
    "react/jsx-no-constructed-context-values": "error",
    "react/no-unused-prop-types": "error",
    "arrow-body-style": "error",
    "curly": "error",
    "no-only-tests/no-only-tests": "error",
    "@typescript-eslint/default-param-last": "error",
    "@typescript-eslint/dot-notation": ["error", { "allowKeywords": true }],
    "@typescript-eslint/method-signature-style": "error",
    "@typescript-eslint/naming-convention": [
      "error",
      {
        "selector": "variable",
        "format": ["camelCase", "PascalCase", "UPPER_CASE"]
      },
      {
        "selector": "function",
        "format": ["camelCase", "PascalCase"]
      },
      {
        "selector": "typeLike",
        "format": ["PascalCase"]
      }
    ],
    "@typescript-eslint/no-unused-expressions": [
      "error",
      {
        "allowShortCircuit": false,
        "allowTernary": false,
        "allowTaggedTemplates": false
      }
    ],
    "@typescript-eslint/no-redeclare": "error",
    "@typescript-eslint/no-shadow": "error",
    "@typescript-eslint/return-await": ["error", "in-try-catch"],
    "camelcase": "warn",
    "no-else-return": ["error", { "allowElseIf": false }],
    "eqeqeq": ["error", "always", { "null": "ignore" }],
    "react/jsx-curly-brace-presence": [2, { "props": "never", "children": "never" }],
    "object-shorthand": ["error", "always"],
    "no-console": "error",
    "no-param-reassign": [
      "error",
      {
        "props": true,
        "ignorePropertyModificationsFor": ["acc", "e"],
        "ignorePropertyModificationsForRegex": ["^assignable[A-Z]"]
      }
    ],
    "@typescript-eslint/no-base-to-string": "error",
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/interface-name-prefix": "off",
    "@typescript-eslint/no-var-requires": "off",
    "@typescript-eslint/no-empty-function": "error",
    "@typescript-eslint/no-inferrable-types": "error",
    "@typescript-eslint/no-unused-vars": "error",
    "@typescript-eslint/explicit-module-boundary-types": "error",
    "react/self-closing-comp": "error",
    "@typescript-eslint/no-unnecessary-condition": "error",
    "react-hooks/exhaustive-deps": "error",
    "prefer-destructuring": [
      "error",
      {
        "VariableDeclarator": {
          "array": false,
          "object": true
        },
        "AssignmentExpression": {
          "array": true,
          "object": false
        }
      },
      {
        "enforceForRenamedProperties": false
      }
    ],
    "react-hooks/rules-of-hooks": "error",
    "import/extensions": "off",
    "import/no-unresolved": "off",
    "import/order": [
      "error",
      {
        "pathGroups": [
          {
            "pattern": "~/**",
            "group": "external",
            "position": "after"
          }
        ],
        "pathGroupsExcludedImportTypes": ["builtin"],
        "groups": [
          "builtin",
          "external",
          "internal",
          "index",
          "sibling",
          "parent",
          "object",
          "unknown"
        ]
      }
    ],
    "no-restricted-properties": [
      "error",
      {
        "property": "sort",
        "message": "Avoid using .sort, use .toSorted instead."
      }
    ],
    "import/newline-after-import": "error",
    "import/no-duplicates": "error",
    "import/no-named-as-default": "error",
    "import/no-extraneous-dependencies": [
      "error",
      {
        "devDependencies": true,
        "optionalDependencies": true
      }
    ],
    "no-relative-import-paths/no-relative-import-paths": [
      "warn",
      {
        "allowSameFolder": true,
        "rootDir": "src",
        "prefix": "~"
      }
    ],
    "prettier/prettier": [
      "error",
      {
        "arrowParens": "always",
        "singleQuote": true,
        "trailingComma": "all",
        "printWidth": 100
      }
    ],
    "react/prop-types": "off",
    "array-callback-return": ["error", { "allowImplicit": true }],
    "prefer-template": "error",
    "no-lone-blocks": "error",
    "no-lonely-if": "error",
    "no-promise-executor-return": "error",
    "no-restricted-globals": [
      "error",
      {
        "name": "isFinite",
        "message": "Use Number.isFinite instead https://github.com/airbnb/javascript#standard-library--isfinite"
      },
      {
        "name": "isNaN",
        "message": "Use Number.isNaN instead https://github.com/airbnb/javascript#standard-library--isnan"
      }
    ],
    "no-sequences": "error",
    "no-undef-init": "error",
    "no-unneeded-ternary": ["error", { "defaultAssignment": false }],
    "no-useless-computed-key": "error",
    "no-useless-return": "error",
    "symbol-description": "error",
    "yoda": "error",
    "func-names": "warn"
  },
  "overrides": [
    {
      "files": ["./src/api/**"],
      "rules": {
        "no-restricted-imports": [
          "off",
          {
            "patterns": ["~/api/**"]
          }
        ]
      }
    },
    {
      "files": ["./src/__tests__/cypress/**/*.ts"],
      "parserOptions": {
        "project": ["./src/__tests__/cypress/tsconfig.json"]
      },
      "extends": [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:prettier/recommended",
        "prettier",
        "plugin:cypress/recommended"
      ]
    },
    {
      "files": ["*.ts", "*.tsx"],
      "excludedFiles": ["**/__mocks__/**", "**/__tests__/**"],
      "rules": {
        "@typescript-eslint/consistent-type-assertions": [
          "error",
          {
            "assertionStyle": "never"
          }
        ]
      }
    },
    {
      "files": ["src/__tests__/cypress/**"],
      "rules": {
        "@typescript-eslint/consistent-type-imports": "error",
        "no-restricted-imports": [
          "error",
          {
            "patterns": [
              {
                "group": [
                  "@patternfly/**"
                ],
                "message": "Cypress tests should only import mocks and types from outside the Cypress test directory."
              }
            ]
          }
        ]
      }
    }
  ]
}
