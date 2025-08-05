export namespace api {
	
	export class RecordsContainer {
	    Records: any;
	    Total: number;
	
	    static createFrom(source: any = {}) {
	        return new RecordsContainer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Records = source["Records"];
	        this.Total = source["Total"];
	    }
	}

}

export namespace manager {
	
	export class InstanceInfo {
	    Module?: any;
	    // Go type: time
	    Started: any;
	
	    static createFrom(source: any = {}) {
	        return new InstanceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Module = source["Module"];
	        this.Started = this.convertValues(source["Started"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class tunnelsList {
	
	
	    static createFrom(source: any = {}) {
	        return new tunnelsList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Manager {
	    Instances: Record<string, InstanceInfo>;
	    // Go type: tunnelsList
	    Tunnels: any;
	
	    static createFrom(source: any = {}) {
	        return new Manager(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Instances = this.convertValues(source["Instances"], InstanceInfo, true);
	        this.Tunnels = this.convertValues(source["Tunnels"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class ChangeServiceState {
	    hash: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangeServiceState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.state = source["state"];
	    }
	}
	export class ServiceHash {
	    hash: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceHash(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	    }
	}
	export class WebService {
	    hash: string;
	    serviceName: string;
	    moduleName: string;
	    listenIP: string;
	    listenPort: number;
	    moduleSettings: any;
	    autoStart: boolean;
	    active: boolean;
	    baseProto: string;
	
	    static createFrom(source: any = {}) {
	        return new WebService(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.serviceName = source["serviceName"];
	        this.moduleName = source["moduleName"];
	        this.listenIP = source["listenIP"];
	        this.listenPort = source["listenPort"];
	        this.moduleSettings = source["moduleSettings"];
	        this.autoStart = source["autoStart"];
	        this.active = source["active"];
	        this.baseProto = source["baseProto"];
	    }
	}

}

