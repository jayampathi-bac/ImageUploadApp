
// set image
$("#btnSrchImg").click(function () {
    let id = $("#imgid").val();
    $.ajax({
        method: "GET",
        url: "http://localhost:8000/load/"+id,
        success: function (res) {


            $("#h1").val(res.id);
            $("#printpat").val(res.imgpath);
            console.log(res.imgpath)

            // $('#imgsrc').attr('src','/ImgUploadApp/temp-img/oneplus.jpg');
            $('#imgsrc').attr('src',res.imgpath);

        },
        error: function (ob, txtStatus, error) {
            console.log(error);
            console.log(txtStatus);
            console.log(ob);
        }
    });
});