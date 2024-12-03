# Pre-requisites for executable documentation

## Summary

When users have long-form documentation that they want to make executable, it is
very common for that documentation to have pre-requisite documentation that needs
to be followed prior to completing the current document. This is especially true
for documentation that is meant to be used as a tutorial or a guide, where the results
from one document are used in another. To bring this functionality to executable
documentation, we nee

To address this issue, we propose the introduction of pre-requisites for executable
documentation that will allow the documentation authors to specify the pre-requisite
documents that are required for their document to be executed successfully.

## Requirements

- [ ] Documentation authors can specify pre-requisites for a document in a new
      section within their document (eg. `## Pre-requisites`)
  - [ ] Pre-requisites items (Other documents) are specified as a list of
        documents that the innovation engine can execute before executing other documents
  - [ ] Loading documentation from markdown source works
  - [ ] Loading documentation from the `source` query parameter works
- [ ] Documentation authors can also specify a validation section that can be
      used to determine if a pre-requisite document has already been executed
      for a given user.
  - [ ] Variables exported from the validation section are carried over into the
        execution of the current document, allowing for the documents to be

## Technical specifications

- If a pre-requisite section is specified within a document, the author must
  provide a validation document, otherwise
  the entire pre-requisite document will be executed.
- Pre-requisites will only be allowed to be nested one level deep, but will
  support deeper levels of nesting in the future.
- Documents will not be allowed to have circular dependencies and these will be
  verified before any execution begins.
- For continuinity between reading executable documents and pre-requisite documents,
  we will offer two ways for the pre-requisite source to be acquired:
  1. The link to the markdown source itself is provided in the pre-requisite
     section (Could be an http path, relative file system path, etc).
  2. The link to the markdown source is provided as a query parameter in the URL
     used to link to the rendered version of the document (I.E. on MS learn,
     github, etc) and must be in the format `?source=<path>`.
