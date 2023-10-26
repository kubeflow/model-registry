name: (Internal) Tracker Template
description: Intended to help with a template for tracking larger grouped items.
title: "[Tracker]: "
labels: ["tracker"]
body:
  - type: textarea
    id: description
    attributes:
      label: Description
      description: A introductory description of the larger task
    validations:
      required:
        true
  - type: input
    id: branch
    attributes:
      label: Target Branch
      description: What is the feature branch to contain this effort? If not known at this time, replace with `TBD`
      placeholder: f/
    validations:
      required: true
  - type: textarea
    id: requirements
    attributes:
      label: Requirements
      description: A series of requirements to consider this tracker complete.
      placeholder: |
        * P0: Show something
        * P2: Allow users to change permissions
    validations:
      required: true
  - type: textarea
    id: ux-issues
    attributes:
      label: Itemized UX Issues
      description: |
        List the tickets that UX will work on.

        Tip: Using a bullet list will help display links to other tickets by unraveling the name and status of that ticket.
      placeholder: |
        * #1234
        * Design mocks - Ticket TBD
    validations:
      required: true
  - type: textarea
    id: dev-issues
    attributes:
      label: Itemized Dev Issues
      description: |
        List the tickets that Development will work on. If unknown at this time, add `TBD`

        Tip: Using a bullet list will help display links to other tickets by unraveling the name and status of that ticket.
      placeholder: |
        * #1234
        * Implement Table Page - Ticket TBD
    validations:
      required: true
  - type: textarea
    id: artifacts
    attributes:
      label: Related artifacts
      description: Any additional artifacts that will help with the tracker goals
    validations:
      required: false