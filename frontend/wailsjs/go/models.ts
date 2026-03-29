export namespace main {
	
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

