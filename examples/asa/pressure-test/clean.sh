read -p "please enter your namespace prefix: " namespace_prefix

read -e -p "please enter start index: " -i 1 start
read -e -p "please enter end index: " -i 5 end
for((i=start;i<=end;i++));  
do   
  currentNS="${namespace_prefix}-${i}"
  kubectl delete App -l carto.run/template-kind=Pressure -n $currentNS

  items=$(kubectl get pods,App -n $currentNS --output 'jsonpath={.items}')
  while [ "$items" != "[]" ] ;do
    echo 'There is alive App, wait for them to be deleted'
    sleep 5
    items=$(kubectl get pods,App -n $currentNS --output 'jsonpath={.items}')
  done
  echo 'Pod/App is cleared, now delete the namespace'
  kubectl delete namespace $currentNS
done
 