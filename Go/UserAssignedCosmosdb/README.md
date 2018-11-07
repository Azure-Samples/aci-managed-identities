# Access CosmosDB from Azure Container Instances Using Manged Identites

This example walks through setting up a User Assigned Managed Identity to get a CosmosDB connection string from an Azure Key Vault.

## Getting Started

To deploy the container, open either a terminal with the Azure CLI installed or check out [CloudShell](https://shell.azure.com/) for a ready-to-go option

Create a resource group

    az group create --name myResourceGroup --location westus

### Create the User Assigned Identity

Next we create a user assigned identity. Once we create and add the permissions, we will be able to use this for multiple container groups.

    az identity create -g myResourceGroup --name myUserIdentity

The output should look close to the following

    {
        "clientId": "0f73d23a-9097-40e7-8753-89d35ced2ff4",
        "clientSecretUrl": "theSecretURL",
        "id": "<longResourceID>",
        "location": "westus",
        "name": "identityName",
        "principalId": "d259206f-b63f-405e-ae52-551b2e769a8b",
        "resourceGroup": "myResourceGroup",
        "tags": {},
        "tenantId": "4720d658-5531-4d8e-889e-41fad886ab7a",
        "type": "Microsoft.ManagedIdentity/userAssignedIdentities"
    }

To make our lives easier, lets add the important information from this output to some environment variables.

    CLIENT_ID=<clientId>
    PRINCIPAL_ID="<principalId>"
    ID="<id>"

Note: if you’re using Powershell, make sure to add the “$” when declaring these variables.

### Create a Key Vault and Set the Permissions

If you don’t already have a Key Vault create, use the following command to create one:

    az keyvault create -g myResourceGroup --name <mykeyvault>

We can quickly add a secret to the new vault

    az keyvault secret set --name SampleSecret --value "MSI Secret!" --vault-name <mykeyvault>

Now, we can give our identity access to the Key Vault

    az keyvault set-policy -n <mykeyvault> --object-id $PRINCIPAL_ID -g myResourceGroup --secret-permissions get

The above command uses the environment variable we set to give our identity “get” permission for secrets in the Key Vault