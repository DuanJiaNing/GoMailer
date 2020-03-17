function registerApp()
{
    echo "Sending request... $1"
    info=`curl -H "Content-Type:application/json" -X POST --data @$2 $1`
    echo $info
}

registerApp "http://8.9.30.183:8080/api/shortcut" $1