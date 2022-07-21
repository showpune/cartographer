#!/bin/bash
set -e

az account set --subscription d51e3ffe-6b84-49cd-b426-0dc4ec660356
sagroup=azdmss-dogfood
saname=mssdevsharedstorage
saassesskey="ymuzl+E/VLStQdBF1i1hXDnQe6AVHsp53lllv+7Wo+BeQ+0NzoWQ5Qvw6Yfr3XGHrIqRf4e0a2eIP9thOJy1wg=="
fileshare=zhiyongtest

# piggy
#filepath="resources/2022030909-fc4994f0-3142-4a5a-a93e-b66e57d433e0"
filepath="directjar/java-native-0.0.1-SNAPSHOT.jar"
# petilic
#filepath="resources/2022030402-b6905895-a8eb-43d6-9e88-cc850be4645a"
#filepath="resources/2022030308-5cb72481-a4c1-4168-b702-55494eb0e9ab"

echo "=========1. select storage account==========="

if [ "$sagroup" = "" ]; then
    echo "which storage account group?"
    read sagroup
fi
if [ "$saname" = "" ]; then
    echo "which storage account?"
    read saname
fi
az configure --defaults group=$sagroup storage=$saname

echo "=========2. set storage account==========="
if [ "$saassesskey" = "" ]; then
    echo "input storage accesskey"
    read saassesskey
fi
az config set storage.account=$saname storage.access_key=$saassesskey

echo "=========3. generate storage file url==========="
if [ "$fileshare" = "" ]; then
    echo "input file share name"
    read fileshare
fi
if [ "$filepath" = "" ]; then
    echo "input path"
    read filepath
fi
fileurl=`az storage file url --share-name $fileshare --path $filepath`

if [ "$fileurl" = "" ]; then
    echo "fileurl is:"
    echo $fileurl
fi
echo "=========4. generate sas==========="
expiredate=`date --date='100 day' +%Y-%m-%dT%TZ`
sas=`az storage file generate-sas --share-name $fileshare --path $filepath  --permissions r --expiry $expiredate`
echo "sas is:"
echo $sas

echo "===========result==========="
echo $fileurl?$sas | tr -d '"'
