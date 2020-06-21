function uploadFile() {
  var files = document.getElementById("file").files;
  var bucket = document.getElementById("bucket").value;
  for(var i=0;i <= files.length-1; i++) {
    var file = files[i]

    var formData = new FormData();

    formData.append("file", file);
    formData.append("bucket", bucket);
    fetch("/upload", {method: "POST", body: formData});

    var fileSize = Math.round((file.size/1000000)*100)/100;

    console.log(`Uploading "${file.name}" with filesize ${fileSize} MB`);

  }
  setTimeout(function (){
    getBucketObjects();
  }, 1000);
}
