import axios from "axios";
import Router from "@/router";
import moment from "moment";

export async function $get(url, options) {
  return new Promise((resolve, reject) => {
    axios
      .get(url, options)
      .then(resp => {
        if (resp.data.error && resp.data.error.match(/auth required/)) {
          if (Router.history.current.name != "Auth") {
            Router.replace({ name: "Auth" });
          }
          return;
        }
        resolve(resp.data);
      })
      .catch(err => {
        reject(err);
      });
  });
}

export async function $post(url, data, options) {
  return new Promise((resolve, reject) => {
    axios
      .post(url, data, options)
      .then(resp => {
        resolve(resp.data);
      })
      .catch(err => {
        if (err.message.match(/401/)) {
          if (Router.history.current.name != "Auth")
            Router.replace({ name: "Auth" });
          return;
        }
        reject(err);
      });
  });
}

export function timeAgo(date) {
  return moment(date, "YYYY-MM-DD HH:mm:ss ZZ z").fromNow();
}

export function randomString(size) {
  return [...Array(size)]
    .map(() => Math.floor(Math.random() * 16).toString(16))
    .join("");
}

export async function showError(ctx, title, message, closable, timer) {
  ctx.$notify({
    title: title,
    message: message,
    type: "error"
  });
}

export function humanFileSize(bytes, si = false, dp = 1) {
  const thresh = si ? 1000 : 1024;
  if (Math.abs(bytes) < thresh) {
    return bytes + " B";
  }

  const units = si
    ? ["kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"]
    : ["KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"];
  let u = -1;
  const r = 10 ** dp;

  do {
    bytes /= thresh;
    ++u;
  } while (
    Math.round(Math.abs(bytes) * r) / r >= thresh &&
    u < units.length - 1
  );

  return bytes.toFixed(dp) + " " + units[u];
}

export function objectToString(o) {
  return o
    ? Object.keys(o)
        .reduce(function(a, k) {
          if (o[k]) {
            if (typeof o[k] == "string") {
              a.push(k + "=" + encodeURIComponent(o[k]));
            }
            if (typeof o[k] == "object") {
              for (let item of o[k]) {
                a.push(k + "[]" + "=" + encodeURIComponent(item));
              }
            }
          }
          return a;
        }, [])
        .join("&")
    : "";
}

export function generateAndDownload(filename, text) {
  var element = document.createElement("a");
  element.setAttribute(
    "href",
    "data:text/plain;charset=utf-8," + encodeURIComponent(text)
  );
  element.setAttribute("download", filename);
  element.style.display = "none";
  document.body.appendChild(element);
  element.click();
  document.body.removeChild(element);
}

export function capitalize(s) {
  if (typeof s !== "string") return "";
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export function addQueryParam(ctx, key, value) {
  let obj = Object.assign({}, ctx.$route.query);
  // ignore eq fields
  if (key in obj) {
    if (obj[key] === value) {
      return;
    } else {
      obj[key] = value;
    }
  } else {
    obj[key] = value;
  }
  ctx.$router.replace({
    ...ctx.$router.currentRoute,
    query: obj
  });
}

export function removeQueryParam(ctx, key) {
  let obj = Object.assign({}, ctx.$route.query);
  if (key in obj) {
    delete obj[key];
    ctx.$router.replace({
      ...ctx.$router.currentRoute,
      query: obj
    });
  }
}

export default {
  humanFileSize: humanFileSize
};
