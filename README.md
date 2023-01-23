# Serverless Gin API

This was a brief two-day project to quickly try to write a CI/CD pipeline for a serverless RESTful API in Go, using both the Gin package and the AWS Labs' proxy integrations for Lambda. It creates routes for the `GET` method, allowing users to either retrieve a fictional company from the sample data provided or else list all of the fictitious corporations available. It also enables basic authorisation for administrators to use `POST` to list new entities, such as Cyberdyne Systems or Wayne Enterprises.

URLs should resemble the following: "https://xxxxxxxxx.execute-api.eu-west-1.amazonaws.com/".

1. GET method
```bash
curl "$URL/v1/companies/1" \
    --header "Content-Type: application/json" \
    --request "GET" \
```
```json
{
    "id":    "1",
    "name":  "Weyland-Yutani Corp",
    "media": "Alien",
    "type":  "Transnational conglomerate",
    "year":  "1979",
}
```
```bash
curl "$URL/v1/companies/" \
    --include \
    --header "Content-Type: application/json" \
    --request "GET" \
```

2. POST method
```bash
curl "$URL/v1/companies?user=admin&password=password" \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{
        "id":    "5",
        "name":  "Cyberdyne Systems",
        "media": "The Terminator",
        "type":  "Defence systems",
        "year":  "1984",
    }'
```

Its runtime context is provided by `cdk.json`, and the pipeline's deployment to AWS itself would be achieved with the CDK command `cdk deploy`. Some features that could be researched and integrated later with a Gin API include support for DynamoDB, Aurora Serverless, and Secrets Manager, to name but a few. Future work might also involve adapting it to deploy via other CI/CD platforms as well; the AWS SAM CLI does accommodate GitLab CI and can work alongside the CDK.
