import axios from "axios";
import Router from "@/router";

async function checkOnServer(token) {
  let r = await axios.get(`/api/sessions/`, {
    headers: { Authorization: "Bearer " + token }
  });

  if (r && r.data && r.data.status == "ok") {
    return true;
  }

  return false;
}

async function setupAxios(token) {
  axios.defaults.headers.common["Authorization"] = "Bearer " + token;
}

// axios.defaults.headers.common['Authorization'] = token
export async function checkAuth(authPageToken) {
  return new Promise(async (resolve, reject) => {
    if (authPageToken) {
      console.log("check authPageToken", authPageToken);
      let ok = await checkOnServer(authPageToken);
      if (ok) {
        localStorage.setItem("token", authPageToken);
        setupAxios(authPageToken);
        Router.push({ name: "Interactions" });
        resolve(true);
        return true;
      }
    } else {
      let localToken = localStorage.getItem("token");
      if (localToken) {
        let ok = await checkOnServer(localToken);
        if (ok) {
          setupAxios(localToken);
          resolve(true);
          return true;
        } else {
          // remove invalid token
          console.log("remove auth token");
          localStorage.removeItem("token");
          resolve(false);
        }
      }
    }
    resolve(false);
  });
}

export function logout() {
  localStorage.removeItem("token");
  Router.push({ name: "Auth" });
}
