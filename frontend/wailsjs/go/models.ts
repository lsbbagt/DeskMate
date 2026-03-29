export namespace main {
	
	export class CodeFile {
	    id: string;
	    name: string;
	    path: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new CodeFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.type = source["type"];
	    }
	}
	export class CodeFolder {
	    id: string;
	    name: string;
	    files: CodeFile[];
	
	    static createFrom(source: any = {}) {
	        return new CodeFolder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.files = this.convertValues(source["files"], CodeFile);
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
	export class Todo {
	    id: string;
	    name: string;
	    description: string;
	    deadline: string;
	    completed: boolean;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Todo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.deadline = source["deadline"];
	        this.completed = source["completed"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

