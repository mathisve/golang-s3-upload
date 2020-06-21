function getBucketObjects() {
  var bucket = document.getElementById("bucket").value;

  var formData = new FormData();
  formData.append("bucket", bucket);

  fetch("/getBucketObjects", {method: "POST", body: formData})
    .then(response => response.json())
    .then(function (data) {
      var table = document.getElementById("table");

      var child = table.firstElementChild;
      while (child) {
  			table.removeChild(child);
  			child = table.firstElementChild;
		  }

      for(var i=0; i < data.length; i++){
          var node = document.createElement("tr");

          var name = document.createElement("td");
          name.innerHTML = data[i].Name;
          node.appendChild(name);

          var size = document.createElement("td");
          size.innerHTML = data[i].Size;
          node.appendChild(size);

          table.appendChild(node);

      }
    });
}
