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

