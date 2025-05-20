<!-- omit in toc -->
# Contributing to InnovationEngine

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles contributions. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions. 🎉

> And if you like the project, but just don't have time to contribute, that's fine. There are other easy ways to support the project and show your appreciation, which we would also be very happy about:
> - Star the project
> - Tweet about it
> - Refer this project in your project's readme
> - Mention the project at local meetups and tell your friends/colleagues

<!-- omit in toc -->
## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Microsoft Open Source Contribution Guide](#microsoft-open-source-contribution-guide)
- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Your First Code Contribution](#your-first-code-contribution)
- [Reviewing Copilot PRs](#reviewing-copilot-prs)
<!--TODO [Improving The Documentation](#improving-the-documentation)-->
- [Styleguides](#styleguides)


## Code of Conduct

This project and everyone participating in it is governed by the
[InnovationEngine Code of Conduct](https://github.com/Azure/InnovationEngine/blob/main/CODE_OF_CONDUCT.md).
By participating, you are expected to uphold this code. Please report unacceptable behavior
to mbifeld@microsoft.com.

## Microsoft Open Source Contribution Guide

This is a Microsoft Open Source project. Please reference to the [Microsoft Open Source Contribtution Guide](https://docs.opensource.microsoft.com/contributing/) for FAQs and general information on contributing to Microsoft Open Source.

## I Have a Question

<!-- If you want to ask a question, we assume that you have read the available [Documentation](TODO: Add Documentation folder).-->
Before you ask a question, it is best to search for existing [Issues](https://github.com/Azure/InnovationEngine/issues) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open an [Issue](https://github.com/Azure/InnovationEngine/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (golang version, operating system, etc), depending on what seems relevant.

We will then address the issue as soon as possible.

## I Want To Contribute

> ### Legal Notice <!-- omit in toc -->
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.

### Reporting Bugs

<!-- omit in toc -->
#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you down for more information. Therefore, we ask you to investigate carefully, collect information, and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the [documentation](./README.md). If you are looking for support, you might want to check [I Have A Question](#i-have-a-question)).
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in [Issues](https://github.com/Azure/InnovationEngine/issues).
- Also make sure to search the internet (including Stack Overflow) to see if users outside of the GitHub community have discussed the issue.
- Collect information about the bug:
- Stack trace (Traceback)
- OS, Platform and Version (Windows, Linux, macOS, x86, ARM, etc)
- Version of the golang, make, etc depending on what seems relevant.
- Possibly your input and the output
- Can you reliably reproduce the issue? And can you also reproduce it with older versions?

<!-- omit in toc -->
#### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities, or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead, sensitive bugs must be sent by email to mbifeld@microsoft.com.
<!-- You may add a PGP key to allow the messages to be sent encrypted as well. -->

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

- Open an [Issue](https://github.com/Azure/InnovationEngine/issues/new). (Since we can't be sure at this point whether it is a bug or not, we ask you not to talk about a bug yet and not to label the issue.)
- Explain the behavior you would expect and the actual behavior.
- Please provide as much context as possible and describe the *reproduction steps* that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.
- Provide the information you collected in the previous section.

Once it's filed:

- The project team will label the issue accordingly.
- A team member will try to reproduce the issue with your provided steps. If there are no reproduction steps or no obvious way to reproduce the issue, the team will ask you for those steps and mark the issue as `needs-repro`. Bugs with the `needs-repro` tag will not be addressed until they are reproduced.
- If the team is able to reproduce the issue, it will be marked `needs-fix`, as well as possibly other tags (such as `critical`), and the issue will be left to be [implemented by someone](#your-first-code-contribution).

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for InnovationEngine, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

<!-- omit in toc -->
#### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the [documentation](./README.md) carefully and find out if the functionality is already covered, maybe by an individual configuration.
- Perform a [search](https://github.com/Azure/InnovationEngine/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.
- Find out whether your idea fits within the scope and aims of the project. It's up to you to make a strong case to convince the project's developers of the merits of this feature. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset. If you're just targeting a minority of users, consider writing an add-on/plugin library.

<!-- omit in toc -->
#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub Issues](https://github.com/Azure/InnovationEngine/issues).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
- You may want to **include screenshots and animated GIFs** which help you demonstrate the steps or point out the part which the suggestion is related to.
- **Explain why this enhancement would be useful** to most InnovationEngine users. You may also want to point out the other projects that solved it better and which could serve as inspiration.

### Your First Code Contribution
#### Innovation Engine
To get started with developing features for the Innovation Engine itself, you 
will need `make` & `go`. Once you have those installed and the project cloned 
to a local repository, you can attempt to build the project using:

```bash
make build-all
```

If the build completes, you should be able to start adding features/fixes
to the Innovation Engine codebase. Once you've added new changes, you can test
for regressions using:

```bash
make test-all
```

If implementing a new feature, it is expected to add & update any necessary 
tests for the changes introduced by the feature.

If you're still looking for more information about how to build & run Innovation Engine, 
[README](./README.md) has a more comprehensive guide for how to get started with project
development.

#### Innovation Engine markdown scenarios

If you are contributing to one of the markdown scenarios (executable documents) 
for Innovation Engine, you are expected to follow the installation steps before 
updating/adding your document. This is needed because once you've made changes 
or have added a new scenario, you should test your executable document by 
using the Innovation Engine:

```bash
ie execute <scenario-path>
```

This will attempt to parse your document into an executable scenario, make sure 
that the commands extracted from codeblocks execute successfully, and that 
their corresponding result blocks (if any) also line up with what the command 
returned. Once you get your scenario to execute successfully, you should go ahead 
and make a PR for it!

#### Creating a PR

<!-- TODO: Create a PR template -->
When creating a PR, please include as much context as possible. At minimum, this should include what the PR does and the testing strategies for it.

If your PR is a work in progress, please label it as a draft and include 'WIP' at the beginning of the PR title.

#### Reviewing Copilot PRs

For guidance on how to review pull requests that were authored by GitHub Copilot, please refer to our [Reviewing Copilot PRs](./docs/REVIEWING_COPILOT_PR.md) guide.

<!--### Improving The Documentation
TODO-->

## Styleguides
For working on the Innovation Engine, `go fmt` is what is used to format the 
code for the project. 

The commit style for individual commits doesn't necessarily matter as 
all commits from a PR branch will be squashed and merged into the main 
branch when PRs are completed.

<!-- omit in toc -->
## Attribution
This guide is based on the **contributing-gen**. [Make your own](https://github.com/bttger/contributing-gen)!
