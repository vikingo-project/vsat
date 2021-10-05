<template>
  <div class="container">
    <div class="row">
      <div class="col-md-3 col-lg-2 col-xl-2 navbar-collapse sidebar">
        <div class="sidebar-header">
          <a class="vik-brand" href="#">
            <img src="@/assets/logo.svg" alt />
          </a>
          <div id="dismissSidebar">
            <svg
              viewBox="0 0 24 24"
              width="24"
              height="24"
              xmlns="http://www.w3.org/2000/svg"
              class="SideNav__closeIcon__56eIa"
              stroke="currentColor"
              tabindex="0"
            >
              <path d="M1 1L23 23M1 23L23 1" />
            </svg>
          </div>
        </div>

        <div class="category-list">
          <h2>Filter</h2>
          <el-form>
            <el-form-item label="File name">
              <el-input v-model="filter.file_name"></el-input>
            </el-form-item>
            <el-form-item label="Content type">
              <el-select v-model="filter.file_type" clearable>
                <el-option v-for="type in types" :key="type" :value="type"
                  >{{ type }}
                </el-option>
              </el-select>
            </el-form-item>
          </el-form>
          <el-button
            :loading="false"
            class="btn btn-primary btn-lg"
            @click="applyFilter"
            icon="vik vik-search"
          >
            Apply
          </el-button>
        </div>
      </div>
      <div class="col">
        <div class="content" v-loading="filesLoading">
          <div class="market-page">
            <h1>Files</h1>
          </div>
          <div class="row item-container attach-list">
            <div class="col-md-6">
              <div class="h-100 pb-20">
                <article
                  @click="uploadDialogVisible = true"
                  class="btn btn-primary btn-upload"
                >
                  <i class="vik vik-plus"></i>Upload
                </article>
              </div>
            </div>

            <div class="col-md-6" v-for="file in files" :key="file.hash">
              <article class="attach-card">
                <div class="attach-card-left-side">
                  <div class="attach-article-card-title">
                    <h3
                      class="d-flex align-items-center"
                      style="cursor: pointer"
                      @click="downloadFile(file.hash, file.file_name)"
                    >
                      <i class="el-icon-paperclip"></i> {{ file.file_name }}
                    </h3>
                    <span class="attach-source">{{ file.content_type }}</span>
                    <p class="attach-info">
                      {{ fileSize(file.size) }}, {{ timeAgo(file.date) }}
                    </p>
                  </div>
                </div>
                <div class="attach-card-right-side attach-pop">
                  <div class="attach-tags">
                    <span
                      class="tag"
                      v-if="file.interaction_hash"
                      @click="goToInteraction(file.interaction_hash)"
                      >go to interaction</span
                    >
                  </div>
                  <div class="btn-group">
                    <el-dropdown trigger="click" @command="handleCommand">
                      <span
                        class="
                          vik-button
                          btn-sm btn-icon btn
                          vik-button--default
                        "
                        ><i class="el-icon-more"></i
                      ></span>
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item
                          icon="vik vik-download"
                          :command="{
                            cmd: 'download_file',
                            file_name: file.file_name,
                            file_hash: file.hash,
                          }"
                          >Download</el-dropdown-item
                        >
                        <el-dropdown-item
                          class="del"
                          icon="vik vik-delete"
                          :command="{
                            cmd: 'remove_file',
                            file_name: file.file_name,
                            file_hash: file.hash,
                          }"
                          >Remove</el-dropdown-item
                        >
                      </el-dropdown-menu>
                    </el-dropdown>
                  </div>
                </div>
              </article>
            </div>
          </div>
        </div>
      </div>
    </div>
    <el-dialog
      title="Upload a file"
      :visible.sync="uploadDialogVisible"
      class="upload-modal"
    >
      <div class="" style="width: 100%; display: inline-flex">
        <div class="border-block">
          <el-upload
            :headers="extraHaders"
            ref="upload"
            :on-success="onFilesUploaded"
            drag
            action="/api/files/upload/"
            :auto-upload="false"
            :file-list="uploadedFiles"
            multiple
          >
            <i class="vik vik-upload"></i>
            <div class="el-upload__text">
              Drop file here or <em>click to choose</em>
            </div>
            <div class="el-upload__tip" slot="tip">
              files with a size less than 256mb
            </div>
          </el-upload>
        </div>
      </div>

      <span slot="footer" class="dialog-footer">
        <div class="btn-line">
          <button class="btn btn-primary" @click="uploadDialogVisible = false">
            Cancel
          </button>
          <el-button
            :loading="false"
            class="btn btn-success"
            @click="upload()"
            icon="vik vik-save"
          >
            Upload
          </el-button>
        </div>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { bus } from "@/bus.js";
import * as utils from "@/utils";
import DateRenderer from "@/components/DateRenderer";

export default {
  components: { DateRenderer },
  data() {
    return {
      form: {},
      uploadDialogVisible: false,
      files: [],
      applyLoading: false,
      uploadedFiles: [],
      menuItems: [],
      types: [],
      filesLoading: false,
      filter: { file_name: "", file_type: "" },
    };
  },
  created() {
    bus.$on("refresh-files", () => this.loadFiles());
  },
  watch: {
    $route(to, from) {
      this.loadFiles();
    },
  },
  async mounted() {
    this.loadFiles();
    this.loadTypes();
  },
  computed: {
    extraHaders() {
      return {
        Authorization: "Bearer " + localStorage.getItem("token"),
      };
    },
  },
  methods: {
    timeAgo(d) {
      return utils.timeAgo(d);
    },
    handleCommand(item) {
      let { cmd, file_name, file_hash } = item;
      if (cmd === "download_file") {
        this.downloadFile(file_hash, file_name);
      }
      if (cmd === "remove_file") {
        this.removeFile(file_hash);
      }
    },
    goToInteraction(hash) {
      this.$router.push({ name: "Iteractions", query: { hash: hash } });
    },
    async applyFilter() {
      this.applyLoading = true;
      await this.loadFiles();
      this.applyLoading = false;
    },

    async loadFiles() {
      this.filesLoading = true;
      let merged = Object.assign(this.$route.query, { page: this.currentPage });
      merged = Object.assign(merged, this.filter);
      utils.$get(`/api/files/?` + utils.objectToString(merged)).then((data) => {
        this.filesLoading = false;
        if (data.status == "ok") {
          this.files = data.files;
          this.total = data.total;
        } else {
          this.filesLoading = false;
          this.$notify.error({
            title: "Error",
            message: "Failed to load files",
          });
        }
      });
    },
    async loadTypes() {
      utils.$get(`/api/files/types/`).then((data) => {
        if (data.status === "ok") {
          this.menuItems = [];
          this.types = data.types;
        }
      });
    },
    downloadFile(hash, name) {
      let anchor = document.createElement("a");
      anchor.download = name;
      anchor.href = `/api/files/download/${hash}/`;
      anchor.click();
    },
    removeFile(hash) {
      let that = this;
      utils
        .$post(`/api/files/remove/`, { hash: hash })
        .then((data) => {
          if (data.status === "error") {
            that.$notify.error({
              title: "Error",
              message: data.error,
            });
            return;
          }
          this.loadFiles();
        })
        .catch((err) => {
          that.$notify.error({
            title: "Error",
            message: err,
          });
        });
    },

    onFilesUploaded() {
      this.uploadDialogVisible = false;
      this.uploadedFiles = [];
      bus.$emit("refresh-files");
    },
    upload() {
      this.$refs.upload.submit();
    },
    fileSize(size) {
      return utils.humanFileSize(size);
    },
    capitalize(s) {
      return utils.capitalize(s);
    },
    menuFilter(items, name) {
      return items.filter((i) => i.name == name);
    },
    isActiveMenuItem(k, v) {
      return this.$route.query[k] && this.$route.query[k] === v;
    },
    addParam(k, v) {
      utils.addQueryParam(this, k, v);
    },
    removeParam(k) {
      utils.removeQueryParam(this, k);
    },
  },
};
</script>