read -p "please enter your namespace prefix: " namespace_prefix
read -e -p "please enter start index: " -i 1 start
read -e -p "please enter end index: " -i 5 end
for((i=start;i<=end;i++));  
do   
  currentNS="${namespace_prefix}-${i}"
  echo "Start to process namespace $currentNS"
  kubectl create secret generic customer-aks --from-file=kubeconfig=./containerapp.yaml -n $currentNS
done