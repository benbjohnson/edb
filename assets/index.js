var width  = 500,
    height = 500;

var padding = 20;
var xAxisHeight = 20;
var brushHeight = 20;
var contextHeight = xAxisHeight + brushHeight;

var layout = d3.layout.pack()
    .size([width - 4, height - contextHeight - padding - 4])
    .value(function(d) { return d.value; });

var svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height)
var g = svg.append("g")
    .attr("transform", "translate(2,2)");

var x = d3.time.scale().range([0, width]);
var xAxis = d3.svg.axis().scale(x).orient("bottom");
var brush = d3.svg.brush()
    .x(x)
    .on("brushend", function() {
        refresh();
    });

var context = g.append("g")
    .attr("class", "context")
    .attr("transform", "translate(0," + (height-contextHeight) + ")");

context.append("g")
  .attr("class", "x axis")

context.append("g")
  .attr("class", "x brush")

var events = [];
d3.json("/events.json", function(error, data) {
    data.forEach(function(d) {
        d.timestamp = d3.time.format.iso.parse(d.timestamp);
    });
    events = data;

    refresh();
});


function refresh() {
    // Extract min/max timestamps and set domain.
    x.domain(d3.extent(events.map(function(d) {
        return d.timestamp;
    })));

    // Constructor actor/repo hierarchy.
    var reposByActor = {id: 0, children:[], lookup: {}};
    var id = 1;
    for(var i=0; i<events.length; i++) {
        var e = events[i];

        // If brush exists, filter by extent.
        if(!brush.empty() && (e.timestamp.valueOf() < brush.extent()[0].valueOf() || e.timestamp.valueOf() > brush.extent()[1].valueOf())) {
            continue;
        }

        // Create actor lookup.
        if(reposByActor.lookup[e.actor] == null) {
            var actor = {id: id++, type: "actor", name: e.actor, children: [], lookup: {}};
            reposByActor.children.push(actor);
            reposByActor.lookup[e.actor] = actor;
        }

        // Create repo under actor.
        var actor = reposByActor.lookup[e.actor];
        if(actor.lookup[e.repository] == null) {
            var repo = {id: id++, type: "repo", name: e.repository, value: 0}
            actor.children.push(repo);
            actor.lookup[e.repository] = repo
        }

        // Increment value of repo.        
        actor.lookup[e.repository].value++;
    }

    g.datum(reposByActor).selectAll(".node").data(layout.nodes, function(d) { return d.id.toString(); })
      .call(function(selection) {
        var enter = selection.enter(),
            exit  = selection.exit();

        selection.select("circle")
          .transition().duration(500)
            .attr("r", function(d) { return d.r; });

        // Create circles for actors and repos.
        nodes = enter.append("g")
          .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"})
        nodes.append("circle")
          .transition().delay(500)
            .attr("r", function(d) { return d.r; });
        nodes.append("title")
        nodes.filter(function(d) { return d.type == "repo"; }).append("text")
          .attr("fill-opacity", 0)
          .transition().delay(500)
            .attr("fill-opacity", 1);

        // Remove exited nodes.
        exit.remove();

        // Update the position of each node.
        selection.attr("class", function(d) { 
            switch(d.type) {
            case "repo": return "repo node";
            case "actor": return "actor node";
            default: return "node";
            }
           })
        selection.transition().duration(500)
          .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"})
        selection.select("title").text(function(d) {
            if(d.type == null) {
                return "";
            }
            return d.name + " - " + d.value;
        });

        selection.select("text")
            .attr("dy", ".3em")
            .style("text-anchor", "middle")
            .text(function(d) {
                return d.name.slice(d.name.indexOf("/")+1).substring(0, d.r / 3);
            });
    })

    // Create and update content.
    // context.append("path")
    //   .attr("class", "area")
    //   .attr("d", area2)
    //   .datum(data);

    context.select(".x.axis")
      .attr("transform", "translate(0," + brushHeight + ")")
      .call(xAxis);

    context.select(".x.brush")
        .call(brush)
        .selectAll("rect")
          .attr("y", -6)
          .attr("height", brushHeight+7);
}
