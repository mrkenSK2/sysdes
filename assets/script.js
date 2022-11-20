/* placeholder file for JavaScript */
const confirm_delete = (id) => {
    if(window.confirm(`Task ${id} を削除します．よろしいですか？`)) {
        location.href = `/task/delete/${id}`;
    }
}
 
const confirm_update = (id) => {
    // 練習問題 7-2
    if(window.confirm(`Task ${id} を更新します．よろしいですか？`)) {
        return true;
    }else{
        return false;
    }
}

function pushEyeButton() {
    var passwdElm = document.getElementById("password");
    var re_enterPasswdElm = document.getElementById("re_enter_password");
    var eyeBtn = document.getElementById("eyeButton");
    if (passwdElm.type === "text") {
        passwdElm.type = "password";
        re_enterPasswdElm.type = "password";
        eyeBtn.className = "fa fa-eye";
    } else {
        passwdElm.type = "text";
        re_enterPasswdElm.type = "text";
        eyeBtn.className = "fa fa-eye-slash";
    }
}

const confirm_userUpdate = () => {
    if(window.confirm(`ユーザ情報を更新します．よろしいですか？`)) {
        return true;
    }else{
        return false;
    }
}

const confirm_withdrawal = () => {
    if(window.confirm(`ユーザ情報を削除します．よろしいですか？`)) {
        location.href = `/user/delete`;
    }
}