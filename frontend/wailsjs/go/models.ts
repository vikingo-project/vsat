export namespace models {
	
	export class WebService {
	    serviceName: string;
	    moduleName: string;
	    listenIP: string;
	    listenPort: number;
	    moduleSettings: string;
	    autoStart: boolean;
	    active: boolean;
	    baseProto: string;
	
	    static createFrom(source: any = {}) {
	        return new WebService(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
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

