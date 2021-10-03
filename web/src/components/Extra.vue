<template>
  <div>
    <ul class="resource-list" v-if="loaded">
      <li :key="event.hash" v-for="event in parsedEvents" class="resource-item">
        <div class="resourse-step-item mb-2">
          <div class="step-item">
            <span class="step-event" v-if="event.name">{{ event.name }}</span>
            <i>
              <DateRenderer :date="event.date" :full="true"></DateRenderer>
            </i>
          </div>
        </div>

        <div
          class="resourse-cont"
          data-old-padding-top=""
          data-old-padding-bottom=""
          style=""
          data-old-overflow=""
        >
          <div
            class="resource-wrap col-auto"
            v-for="(t, i) in event.data"
            :key="i"
          >
            <div class="resource-block">
              <h3>{{ formatName(t.name) }}</h3>
              <div>
                <div v-if="getDataFormat(t.name) == 'hex'">
                  <Hex :content="t.value"></Hex>
                </div>
                <div v-if="getDataFormat(t.name) == 'file'">
                  <router-link
                    :to="{ name: 'Files', query: { hash: t.value } }"
                    >{{ t.value }}</router-link
                  >
                </div>
                <div v-if="getDataFormat(t.name) == 'text'">
                  <pre>{{ replaceNonASCII(t.value) }}</pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
    </ul>
    <div v-else>loading...</div>
  </div>
</template>

<script>
import DateRenderer from "@/components/DateRenderer";
import Hex from "@/components/Hex";

export default {
  components: { DateRenderer, Hex },
  props: ["event", "events"],
  data() {
    return {
      loaded: false,
      full: {},
      parsedEvents: [],
    };
  },
  watch: {
    events() {
      if (this.events) {
        for (let event of this.events) {
          try {
            let data = JSON.parse(event.data);
            let keys = Object.keys(data);
            event.data = [];
            for (let key of keys) {
              event.data.push({ name: key, value: data[key] });
            }
            this.parsedEvents.push(event);
          } catch (e) {}
        }
        this.loaded = true;
      }
      this.loaded = true;
    },
  },
  mounted() {
    if (this.events) {
      this.loaded = true;
    }
  },
  methods: {
    formatName(n) {
      let p = n.split(":");
      if (p.length > 1) {
        return p[1];
      }
      return p[0];
    },
    getDataFormat(n) {
      let p = n.split(":");
      if (p.length > 1) {
        return p[0];
      }
      return "text";
    },
    replaceNonASCII(v) {
      return v.replace(/[^\x20-\x7E^\n^\r^\tab]+/g, ".");
    },
    value(o, k) {
      return o[k];
    },
    keys(o) {
      if (o) {
        return Object.keys(o);
      } else {
        return [];
      }
    },
  },
};
</script> 