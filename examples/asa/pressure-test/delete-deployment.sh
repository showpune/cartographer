read -e -p "please enter your namespace prefix: " -i "test" namespace_prefix

read -p "please enter start index: " start
read -p "please enter end index: " end

for((i=start;i<=end;i++));  
do   
  currentNS="${namespace_prefix}-${i}"
  kubectl delete App -l carto.run/template-kind=Pressure -n $currentNS --wait=false
  items=$(kubectl get App -n $currentNS --output 'jsonpath={.items}')
  while [ "$items" != "[]" ] ;do
    echo 'There is alive App, wait for them to be deleted'
    sleep 5
    items=$(kubectl get App -n $currentNS --output 'jsonpath={.items}')
  done
  echo 'Pod/App is cleared, now delete the namespace'
    kubectl delete deployment -l kapp.k14s.io/app=kapp -n $currentNS --wait=false
done
 