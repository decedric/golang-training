test=$(curl -X POST -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/100) 
respo=$(jq -sr '.[1].code' <<< "$test")
want="201"
echo "$test"
echo "$respo"
echo "$want"
if [ "201" = "$(jq -sr '.[1].code' <<< "$test")" ]; 
then 
  echo "done" 
else 
  echo "no" 
fi


strval1="Ubuntu"
strval2="Windows"

#Check equality two string variables

if [ $strval1 == $strval2 ]; then
  echo "Strings are equal"
else
  echo "Strings are not equal"
fi

#Check equality of a variable with a string value

if [ $strval1 == "Ubuntu" ]; then
  echo "Linux operating system"
else
  echo "Windows operating system"
fi

