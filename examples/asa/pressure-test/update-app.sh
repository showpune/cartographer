read -p "please enter your namespace prefix: " namespace_prefix

read -e -p "please enter start index: " -i 1 start
read -e -p "please enter end index: " -i 5 end
read -e -p "please enter App Number: " -i 50 number


for((i=start;i<=end;i++));  
do   
  currentNS="${namespace_prefix}-${i}"
  folder="${currentNS}-update" 
  rm -r temp/$folder
  mkdir temp
  mkdir temp/$folder
  cd temp/$folder
  for((j=1;j<=$number;j++));  
  do
    appNumberCode="code-${j}"
    sed s/appNumberCode/$appNumberCode/g ../../app.yaml.template-update.yaml > $appNumberCode.yaml
  done
  sed s/currentNS/$currentNS/g *.yaml > $currentNS.yaml
  kubectl apply -f $currentNS.yaml -n $currentNS
  cd ../..
  rm -rf temp/$folder
done
 