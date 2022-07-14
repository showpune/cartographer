read -p "please enter your namespace prefix: " namespace_prefix

read -e -p "Input the subscription id for container app: " -i "d51e3ffe-6b84-49cd-b426-0dc4ec660356" subscriptionid
read -e -p "Input the resourcegroup for container app: " -i "zhiyongli" resourcegroup
read -e -p "Input the ask name for container app: " -i "containerapp" aksname
az account set --subscription $subscriptionid
az aks get-credentials --resource-group $resourcegroup --name $aksname --admin --overwrite-existing --file ./containerapp.yaml

read -e -p "please enter start index: " -i 1 start
read -e -p "please enter end index: " -i 5 end
read -e -p "please enter App Number: " -i 5 number
for((i=start;i<=end;i++));  
do   
  currentNS="${namespace_prefix}-${i}"
  echo "Start to process namespace $currentNS"
  kubectl create namespace $currentNS
  kubectl create secret generic customer-aks --from-file=kubeconfig=./containerapp.yaml -n $currentNS
  kubectl apply -f ../namespaced/service-account.yaml -n $currentNS

  kubectl create namespace $currentNS --kubeconfig ./containerapp.yaml

  rm -r temp/$currentNS
  mkdir temp
  mkdir temp/$currentNS

  cd temp/$currentNS


  sed s/currentNS/$currentNS/g ../../kapp-child.yaml > 0000.yaml
  for((j=1;j<=$number;j++));  
  do
    appNumberCode="code-${j}"
    sed s/appNumberCode/$appNumberCode/g ../../app.yaml.template.yaml > $appNumberCode.yaml
  done
  sed s/currentNS/$currentNS/g *.yaml > $currentNS.yaml
  kubectl apply -f $currentNS.yaml -n $currentNS
  
  cd ../..
  rm -rf temp/$currentNS
done
 