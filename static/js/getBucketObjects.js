function getBucketObjects() {
  var bucket = document.getElementById("bucket").value;

  var formData = new FormData();
  formData.append("bucket", bucket);

  fetch("/getBucketObjects", {method: "POST", body: formData})
    .then(response => response.json())
    .then(function (data) {
        var items = document.getElementById("itemCount");


        var amountOfItems
        if (data.length >= 1000) {
            amountOfItems = "+1000"
        } else {
            amountOfItems = data.length
        }

        items.innerText = `Items: ${amountOfItems}`;

      var table = document.getElementById("table");

      var child = table.firstElementChild;
      while (child) {
  			table.removeChild(child);
  			child = table.firstElementChild;
      }

        var node = document.createElement("tr");

        var name = document.createElement("td");
        name.innerHTML = "Object";
        name.setAttribute("style", "font-size: 1.3rem;")
        node.appendChild(name);

        var size = document.createElement("td");
        size.innerHTML = "Size in KB";
        size.setAttribute("style", "font-size: 1.3rem;")
        node.appendChild(size);

        table.appendChild(node);


      for(var i=0; i < data.length; i++){
          var node = document.createElement("tr");

          var name = document.createElement("td");
          name.innerHTML = data[i].Name;
          node.appendChild(name);

          var size = document.createElement("td");
          size.innerHTML = Math.round(data[i].Size / 10) / 100 ;
          node.appendChild(size);

          table.appendChild(node);

      }
    });
}
