<template>
  <div>
    <el-form :model="settings" ref="settingsForm" @submit.native.prevent>
      <div class="d-flex">
        <el-form-item
          label="Min passive port"
          class="mr-4"
          prop="minPassivePort"
        >
          <el-input-number
            v-model="settings.minPassivePort"
            suffix-icon="el-icon-date"
            placeholder="Port number"
            :min="1"
            :key="65535"
          ></el-input-number>
        </el-form-item>
        <el-form-item label="Max passive port" prop="maxPassivePort">
          <el-input-number
            v-model="settings.maxPassivePort"
            suffix-icon="el-icon-date"
            placeholder="Port number"
            :min="1"
            :key="65535"
          ></el-input-number>
        </el-form-item>
      </div>

      <el-form-item label="File system" prop="tree">
        <div class="custom-tree-container">
          <div class="block">
            <el-tree
              :data="tree"
              node-key="id"
              default-expand-all
              :expand-on-click-node="false"
            >
              <span class="custom-tree-node" slot-scope="{ node, data }">
                <span v-if="data.folder">{{ node.label }}</span>
                <div class="ml-4" v-if="data.folder">
                  <el-dropdown @command="handleCommand" trigger="click">
                    <span class="vik-button btn-sm btn-icon btn btn-mute"
                      ><i class="el-icon-more"></i
                    ></span>
                    <el-dropdown-menu slot="dropdown">
                      <el-dropdown-item
                        :command="{
                          cmd: 'append_file',
                          node: node,
                          data: data,
                        }"
                        icon="el-icon-document-add"
                        >Add File</el-dropdown-item
                      >
                      <el-dropdown-item
                        :command="{
                          cmd: 'append_folder',
                          node: node,
                          data: data,
                        }"
                        icon="el-icon-folder-add"
                        >Add Folder</el-dropdown-item
                      >
                      <el-dropdown-item
                        :command="{
                          cmd: 'remove',
                          node: node,
                          data: data,
                        }"
                        class="del"
                        icon="vik vik-delete"
                        >Remove</el-dropdown-item
                      >
                    </el-dropdown-menu>
                  </el-dropdown>
                </div>
                <div class="d-flex align-content-center" v-else>
                  <el-select
                    v-model="mappings[data.id].value"
                    placeholder="Choose file"
                    suffix-icon="el-icon-file"
                    filterable
                    @change="setNodeFilename(data)"
                  >
                    <el-option
                      class="file-select"
                      :key="file.hash"
                      v-for="file in files"
                      :label="file.file_name"
                      :value="file.hash"
                    ></el-option>
                  </el-select>
                  <span
                    v-if="mappings[data.id].value"
                    class="mx-3"
                    style="align-self: center; line-height: 1"
                    >as
                    <el-input
                      type="text"
                      v-model="mappings[data.id].asFilename"
                    ></el-input>
                  </span>

                  <div class="btn btn-sm btn-mute" @click="remove(node, data)">
                    <i class="el-icon-close"></i>
                  </div>
                </div>
              </span>
            </el-tree>
          </div>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>
<script>
import * as utils from "@/utils.js";
import Vue from "vue";

var boot = [
  {
    id: 1,
    label: "/",
    folder: true,
    children: [],
  },
];
export default {
  props: ["preset"],
  name: "ftp-settings",
  data() {
    return {
      settings: { fs: {}, minPassivePort: 60000, maxPassivePort: 60100 },
      tree: [],
      files: [],
      newFolderName: "new folder",
      mappings: {},
    };
  },
  async mounted() {
    if (this.preset) {
      this.parsePreset();
    } else {
      this.tree = boot;
    }
    this.files = await this.loadFiles();
  },
  methods: {
    posix(path) {
      if (path.length === 0) return ".";
      if (path === "/") return "/";
      var code = path.charCodeAt(0);
      var hasRoot = code === 47; /*/*/
      var end = -1;
      var matchedSlash = true;
      for (var i = path.length - 1; i >= 1; --i) {
        code = path.charCodeAt(i);
        if (code === 47 /*/*/) {
          if (!matchedSlash) {
            end = i;
            break;
          }
        } else {
          // We saw the first non-path separator
          matchedSlash = false;
        }
      }

      if (end === -1) return hasRoot ? "/" : ".";
      if (hasRoot && end === 1) return "//";
      return path.slice(0, end);
    },
    // TODO: https://gist.github.com/stephanbogner/4b590f992ead470658a5ebf09167b03d
    treeify(files) {
      let that = this;
      files = files.reduce(function (tree, f) {
        var dir = that.posix(f.path);
        console.log("dir", dir, tree);
        if (tree[dir]) {
          tree[dir].children.push(f);
        } else {
          tree[dir] = { implied: true, children: [f] };
        }

        if (tree[f.path]) {
          f.children = tree[f.path].children;
        } else {
          f.children = [];
        }

        return (tree[f.path] = f), tree;
      }, {});

      return Object.keys(files).reduce(function (tree, f) {
        if (files[f].implied) {
          return tree.concat(files[f].children);
        }
        return tree;
      }, []);
    },
    parsePreset() {
      let local = Object.assign({}, this.preset);
      let fs = local.fs;

      let normalized = [];
      let paths = Object.keys(fs);
      for (let path of paths) {
        var id = utils.randomString(6);
        Vue.set(this.mappings, id, {
          filename: "",
          value: "",
          asFilename: "",
        });

        normalized.push({
          id: id + "",
          label: fs[path].name,
          name: fs[path].name,
          folder: fs[path].dir,
          value: "value" in fs[path] ? fs[path].value : "",
          children: [],
          path: "/" + path,
        });
        if ("value" in fs[path]) {
          this.mappings[id].value = fs[path].value;
        }
        if ("name" in fs[path]) {
          this.mappings[id].asFilename = fs[path].name;
        }
      }

      this.tree = this.treeify(normalized);
      this.settings.minPassivePort = local.minPassivePort;
      this.settings.maxPassivePort = local.maxPassivePort;
    },
    handleCommand(item) {
      let { cmd, node, data } = item;
      if (cmd === "append_file") {
        this.appendFile(data);
      }
      if (cmd === "append_folder") {
        this.appendFolder(data);
      }
      if (cmd === "remove") {
        this.remove(node, data);
      }
    },
    getFilenameByHash(hash) {
      for (let f of this.files) {
        if (f.hash === hash) {
          return f.file_name;
        }
      }
      return utils.randomString(5);
    },
    async loadFiles() {
      utils
        .$get(`/api/files/?` + utils.objectToString({ page: this.currentPage }))
        .then((data) => {
          if (data.status == "ok") {
            this.files = data.files;
            this.total = data.total;
          }
        })
        .catch((err) => {
          this.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },

    generateNewFolderName(items, counter = 0) {
      if (items) {
        let hasName = false;
        for (let item of items) {
          console.log("item", item);
          if (
            item.label == this.newFolderName ||
            item.label == `${this.newFolderName} (${counter})`
          ) {
            hasName = true;
            counter++;
            break;
          }
        }
        if (hasName) {
          return this.generateNewFolderName(items, counter);
        }
      }

      return `${this.newFolderName} (${counter})`;
    },
    setNodeFilename(data) {
      let filename = this.getFilenameByHash(this.mappings[data.id].value);
      this.mappings[data.id].asFilename = filename;
    },

    appendFolder(data) {
      let folderName = this.generateNewFolderName(data.children);
      const newChild = {
        id: utils.randomString(6),
        label: folderName,
        children: [],
        folder: true,
      };
      if (!data.children) {
        this.$set(data, "children", []);
      }
      data.children.push(newChild);
    },
    appendFile(data) {
      const newChild = {
        id: utils.randomString(6),
        label: "",
        children: [],
      };
      data.children.push(newChild);

      Vue.set(this.mappings, newChild.id, {
        filename: "",
        value: "",
        asFilename: "",
      });
      // this.mappings[newChild.id] = ;
      console.log("mappings", this.mappings);
    },

    remove(node, data) {
      const parent = node.parent;
      const children = parent.data.children || parent.data;
      const index = children.findIndex((d) => d.id === data.id);
      children.splice(index, 1);
    },

    defaultPort() {
      return 21;
    },
    fsWalk(path, items) {
      for (let item of items) {
        if (item.folder) {
          this.settings.fs[path + item.label] = {
            dir: item.folder ? item.folder : false,
            name: item.label,
          };
        } else {
          this.settings.fs[path + this.mappings[item.id].asFilename] = {
            dir: false,
            name: this.mappings[item.id].asFilename,
            value: this.mappings[item.id].value,
          };
        }

        if (item.children && item.children.length > 0) {
          this.fsWalk(path + item.label, item.children);
        }
      }
    },
    validate() {
      this.settings.fs = {}; // reset
      return new Promise((resolve, reject) => {
        this.fsWalk("", this.tree);
        resolve(this.settings);
      });
    },
  },
};
</script>
<style>
</style>