var width  = 500,
    height = 500;

var layout = d3.layout.pack()
    .size([width - 4, height - 4])
    .value(function(d) { return d.value; });

var svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height)
var g = svg.append("g")
    .attr("transform", "translate(2,2)");

d3.json("http://45.55.179.106:13000/events.json", function(error, events) {
    var reposByActor = {children:[], lookup: {}};

    for(var i=0; i<events.length; i++) {
        var e = events[i];

        // Create actor lookup.
        if(reposByActor.lookup[e.actor] == null) {
            var actor = {type: "actor", name: e.actor, children: [], lookup: {}};
            reposByActor.children.push(actor);
            reposByActor.lookup[e.actor] = actor;
        }

        // Create repo under actor.
        var actor = reposByActor.lookup[e.actor];
        if(actor.lookup[e.repository] == null) {
            var repo = {type: "repo", name: e.repository, value: 0}
            actor.children.push(repo);
            actor.lookup[e.repository] = repo
        }

        // Increment value of repo.        
        actor.value++;
        actor.lookup[e.repository].value++;
    }

    g.datum(reposByActor).selectAll(".node").data(layout.nodes).call(function(selection) {
        // Create circles for actors and repos.
        var enter = selection.enter();
        enter.append("g")
          .attr("class", "node")
          .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"})
          .append("circle");

        // Update the position of each node.
        selection.selectAll("circle").attr("r", function(d) { return d.r; });
        selection.append("title").text(function(d) { return d.name + " - " + d.value; });

        selection.filter(function(d) { return d.type == "repo"; }).append("text")
            .attr("dy", ".3em")
            .style("text-anchor", "middle")
            .text(function(d) {
                return d.name.slice(d.name.indexOf("/")+1).substring(0, d.r / 3);
            });
    })
})
