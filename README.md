# Pull request Querier

Pull request workflow made easy!

## Why

As developers, pull requests sit at the center of our workflow.
Each pull request requires different action depending on its state:

- When someone assigns a pull request to me, I need to review it.
- When someone comments or requests changes on **my** pull request,
  I need to review it and maybe make some changes and re-request reviews.
- When the key reviewers approve my pull request, I need to merge it.

When juggling multiple pull requests, it takes
quite a bit of brain power to remember when we should do what.
**The asynchronous working style makes it even harder.**

## Usage

### Installtion

1. Install by

```
go install github.com/mz026/qp
```

2. Prepare your credential file.
  * Copy the `.qp.credential.sample.yaml` to your home directory and rename it to `.qp.credential.yaml`.
  * [Create a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) on Github.
  * The token should have the following accesses:
    * repo
    * read:org
    * read:repo\_hook
    * read:user
    * read:discussion
  * Fill in the token into the credential file.
  * Fill in the organization where you want to query your pull requests. *You can remove the organization section, and no org filter will be applied.*

### Run!

  ```
  qp
  ```

## LICENCE

MIT
