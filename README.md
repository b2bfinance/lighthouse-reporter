# lighthouse-reporter

Report whole site metrics using Lighthouse by Google.

## Install

```bash
$ go get -v github.com/b2bfinance/lighthouse-reporter/cmd/lhreporter
```

## Usage

```bash
$ lhreport configuration.json
```

Give the output a reference, this will be the file name in Google Cloud Bucket and on other storage just in the file it self.

```bash
$ lhreport configuration.json mySiteResults
```

## Storing the scores

In the configuration there is a `storagePath` key, this can either be any of the following.

- A local file name, this will store the results in the file provided.
- A Google Cloud Storage path for example `gs://myBucketName/path/Prefix` and will go in the path with the following appended to it `<reference>.json`
- A HTTP endpoint where the results will be sent using the `POST` method.

## Testing results

You can use this tool in a CI environment with the configuration specifying a `minimumPageScore` and/or `minimumMeanScore` both of which an object containing `performance`, `accessibility`, `bestPractises` and `seo` scores that must be met.

If these scores are provided a test suite will run against the results and failing to meet the requirements a non zero exit code is returned.

## Endpoints

We require a base URL that is known in the configuration as `remote`, this URL will be base too all found paths. Paths can be loaded from a sitemap by specifying a URL or local file path to an XML sitemap. If you have paths that are not in the sitemap but would need to be checked you can provide these paths in a string array with key `customPaths`

## Example configuration

```json
{
  "remote": "http://localhost:8000",
  "minimumPageScore": {
    "performance": 90,
    "accessibility": 90,
    "bestPractises": 90,
    "seo": 90
  },
  "minimumMeanScore": {
    "performance": 90,
    "accessibility": 90,
    "bestPractises": 90,
    "seo": 90
  },
  "siteMap": "sitemap.xml",
  "customPaths": [
    "/path-not-in-sitemap"
  ],
  "storagePath": "gs://logsBucket/sitename",
  "lighthouseArgs": []
}
```
