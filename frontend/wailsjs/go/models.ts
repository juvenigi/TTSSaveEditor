export namespace properties {
	
	export class ApplicationPropertiesView {
	    OsPathSeparator: string;
	    GameDir: string;
	    PackDataDir: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplicationPropertiesView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OsPathSeparator = source["OsPathSeparator"];
	        this.GameDir = source["GameDir"];
	        this.PackDataDir = source["PackDataDir"];
	    }
	}

}

