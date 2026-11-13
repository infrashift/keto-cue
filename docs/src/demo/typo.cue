package keto

namespaces: {
	User: {}
	Group: relations: members: ["User", "Group#members"]
	Folder: {
		relations: {
			parents: ["Folder"]
			owners: ["User", "Group#members"]
			viewers: ["User", "Group#members"]
		}
		permissions: {
			edit: or: ["owners", "parents->edit"]
			view: or: ["viewers", "edit", "parents->view"]
		}
	}
	Doc: {
		relations: {
			parents: ["Folder"]
			owners: ["User"]
			blocked: ["Usr"]
		}
		permissions: {
			edit: or: ["owners", "parents->edit"]
			view: and: [{or: ["edit", "parents->veiw"]}, {not: "blocked"}]
		}
	}
}
