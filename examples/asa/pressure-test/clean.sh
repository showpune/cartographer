read -p "please enter your namespace prefix: " namespace_prefix

read -e -p "please enter start index: " -i 1 start
read -e -p "please enter end index: " -i 5 end
read -e -p "please enter App Number: " -i 50 number

read -e -p "Delete Namespace: " -i N delete_namespace

for((i=start;i<=end;i++));  
do   
  currentNS="${namespace_prefix}-${i}"
  folder="${currentNS}-delete" 
  rm -r temp/$folder
  mkdir temp
  mkdir temp/$folder
  cd temp/$folder
  for((j=1;j<=$number;j++));  
  do
    appNumberCode="code-${j}"
    sed s/appNumberCode/$appNumberCode/g ../../app.yaml.template-cancel.yaml > $appNumberCode.yaml
  done
  sed s/currentNS/$currentNS/g *.yaml > $currentNS.yaml
# kubectl apply -f $currentNS.yaml -n $currentNS
  cd ../..
  rm -rf temp/$folder

  kubectl delete App -l carto.run/template-kind=Pressure -n $currentNS --wait=false

  kubectl delete deployment -l kapp.k14s.io/app=kapp -n $currentNS --wait=false
  

  if [ $delete_namespace == "Y" ] ;then
  items=$(kubectl get pods,App -n $currentNS --output 'jsonpath={.items}')
  while [ "$items" != "[]" ] ;do
    echo 'There is alive App, wait for them to be deleted'
    sleep 5
    items=$(kubectl get pods,App -n $currentNS --output 'jsonpath={.items}')
  done
  echo 'Pod/App is cleared, now delete the namespace'
  kubectl delete namespace $currentNS
  fi
done
 