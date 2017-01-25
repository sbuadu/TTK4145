mkdir -p driver;
cp ../driver/io.c driver;
cp ../driver/io.h driver;
cp ../driver/channels.h driver;

mkdir -p sim_server;
cp ../simulator/server/* sim_server;
cp ../simulator/simulator.con .;
