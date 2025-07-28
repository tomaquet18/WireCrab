export namespace devices {
	
	export class Device {
	    name: string;
	    description: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.type = source["type"];
	    }
	}

}

export namespace tshark {
	
	export class ProtocolInfo {
	    Name: string;
	    Detail: any;
	    Child?: ProtocolInfo;
	
	    static createFrom(source: any = {}) {
	        return new ProtocolInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Detail = source["Detail"];
	        this.Child = this.convertValues(source["Child"], ProtocolInfo);
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

export namespace types {
	
	export class PacketMeta {
	    Timestamp: string;
	    SrcIP?: string;
	    DstIP?: string;
	    SrcPort?: string;
	    DstPort?: string;
	    Protocol?: string;
	    Length?: number;
	
	    static createFrom(source: any = {}) {
	        return new PacketMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Timestamp = source["Timestamp"];
	        this.SrcIP = source["SrcIP"];
	        this.DstIP = source["DstIP"];
	        this.SrcPort = source["SrcPort"];
	        this.DstPort = source["DstPort"];
	        this.Protocol = source["Protocol"];
	        this.Length = source["Length"];
	    }
	}
	export class CapturedPacket {
	    meta: PacketMeta;
	    parsed?: tshark.ProtocolInfo;
	
	    static createFrom(source: any = {}) {
	        return new CapturedPacket(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], PacketMeta);
	        this.parsed = this.convertValues(source["parsed"], tshark.ProtocolInfo);
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

