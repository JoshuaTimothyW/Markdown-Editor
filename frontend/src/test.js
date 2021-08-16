const axios = require("axios").default

cred = {
    "txt_Username":"Epuser",
    "txt_Password":"userEp",
}

axios.defaults.baseURL = "http://192.168.18.1";

axios.post("/login.cgi",cred)
.then( (res) => {
    console.log(res.data);
})