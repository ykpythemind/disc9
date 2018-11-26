document.addEventListener("DOMContentLoaded", function(event) {
  var inputElm = document.querySelector('input');
  var previewElm = document.querySelector('.preview');
  var stat = document.querySelector('.status');

  var curFiles = []

  var clearPreview = function (prev) {
    while(prev.firstChild) {
      prev.removeChild(prev.firstChild);
    }
  }

  var createItem = function(item, file, i) {
    var image = document.createElement('img');
    image.style = 'width: 100%;'

    // TODO それぞれに再選択できるようにする

    if (!file) {
      image.src = 'http://design-ec.com/d/e_others_51/l_e_others_510.png' // FIXME
      item.appendChild(image)
      return
    }

    var para = document.createElement('p');
    if(validFileType(file)) {
      image.src = window.URL.createObjectURL(file);
      item.appendChild(image);
      para.textContent = '' + ( i + 1 );
    } else {
      para.textContent = 'ファイル ' + file.name + ' はファイル形式が有効ではありません。選択し直してください';
    }
    item.appendChild(para);
  }

  var updateImageDisplay = function () {
    clearPreview(previewElm)

    curFiles = inputElm.files;
    if(curFiles.length === 0) {
      stat.textContent = 'アップロードするファイルが選択されていません';
    } else {
      var list = document.createElement('ol');
      previewElm.appendChild(list);
      for(var i = 0; i < 9; i++) {
        var listItem = document.createElement('li');
        createItem(listItem, curFiles[i], i)
        list.appendChild(listItem);
      }
    }
  }

  var uploadButtonClicked = function(e) {
    e.preventDefault()

    var formData = new FormData();

    for (var i = 0; i < 9; i++) {
      formData.append('files', curFiles[i], curFiles[i].name)
    }

    fetch('/upload', {
      method: 'POST',
      body: formData
    }).then(response => {

      if (response.ok) {

        response.blob().then(blob => {
          var img = document.createElement('img');
          var objectURL = URL.createObjectURL(blob);
          img.src = objectURL;
          clearPreview(previewElm)
          previewElm.appendChild(img)
        });

      } else {
        response.json().then(function(data) {
          stat.textContent = data.message
        });
      }
    });
  }

  var button = document.querySelector('#upload');
  button.addEventListener('click', uploadButtonClicked)

  inputElm.style.opacity = 0;
  inputElm.addEventListener('change', updateImageDisplay);
});


var fileTypes = [
  'image/jpeg',
  'image/pjpeg',
  'image/png'
]

function validFileType(file) {
  for(var i = 0; i < fileTypes.length; i++) {
    if(file.type === fileTypes[i]) {
      return true;
    }
  }
  return false;
}
