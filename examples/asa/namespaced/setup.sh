read -p "please enter your namespace: " namespace
kubectl create namespace $namespace
kubectl config set-context --current --namespace=$namespace
export kubenamespace=$namespace

read -s -p "please enter your private password: " dockerpassword
echo 'Create secret to write the customer image'
kubectl create secret docker-registry private-registry-credentials --docker-server=https://index.docker.io/v1/ --docker-username=showpune --docker-password=$dockerpassword --docker-email=zhiyongli@microsoft.com 

echo 'Create secret to read the cluster resource'
kubectl create secret docker-registry shared-registry-credentials --docker-server=serviceacr.azurecr.io --docker-username=serviceacr --docker-password=1dXPkDXTRv+DlZxi/n8HRC0IjVdKRadY --docker-email=zhiyongli@microsoft.com

read -e -p "Input the subscription id for container app: " -i "d51e3ffe-6b84-49cd-b426-0dc4ec660356" subscriptionid
read -e -p "Input the resourcegroup for container app: " -i "zhiyongli" resourcegroup
read -e -p "Input the ask name for container app: " -i "containerapp" aksname
az account set --subscription $subscriptionid
az aks get-credentials --resource-group $resourcegroup --name $aksname --admin --overwrite-existing --file ./containerapp.yaml

echo 'Create secret to access customer aks'
kubectl create secret generic customer-aks --from-file=kubeconfig=./containerapp.yaml

echo 'Create namespace sa, role and rolebinding'
kubectl apply -f service-account.yaml

