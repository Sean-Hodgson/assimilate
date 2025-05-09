 

assimilateCNCARMExeName = assimCNCARM
assimilateCNCx64ExeName = assimCNCx64

assimilateProxyARMExeName = assimProxARM
assimilateProxyx64ExeName = assimProxx64

assimilateEndpointARMExeName = assimEndARM
assimilateEndpointx64ExeName = assimEndx64

#build for x arch
assimARMBuildDirections = env GOOS=linux GOARCH=arm64
assimx64BuildDirections = env GOOS=linux GOARCH=amd64

MakeBuildDir:
	mkdir -p build

clean:
	rm -rf build

#stuff for the GOCNC
ARMVictim: MakeBuildDir
	cd endpoint && $(assimARMBuildDirections) go build -o $(assimilateEndpointARMExeName)
	cd endpoint && mv $(assimilateEndpointARMExeName) ../build/$(assimilateEndpointARMExeName)

Victim: MakeBuildDir
	cd C2C_Server && $(assimX64BuildDirections) go build -o $(assimilateEndpointx64ExeName)
	cd C2C_Server && mv $(assimilateEndpointx64ExeName) ../build/$(assimilateEndpointx64ExeName)


#stuff for the GOvictim
ARMProx: MakeBuildDir
	cd networking && $(assimARMBuildDirections) go build -o $(assimilateProxyARMExeName)
	cd networking && mv $(assimilateProxyARMExeName) ../build/$(assimilateProxyARMExeName)

Prox: MakeBuildDir
	cd C2C_Server && $(assimX64BuildDirections) go build -o $(assimilateProxyx64ExeName)
	cd C2C_Server && mv $(assimilateProxyx64ExeName) ../build/$(assimilateProxyx64ExeName)


#stuff for the GOPRoxy stuff
ARMC2: MakeBuildDir
	cd C2C_Server && $(assimARMBuildDirections) go build -o $(assimilateCNCARMExeName)
	cd C2C_Server && mv $(assimilateCNCARMExeName) ../build/$(assimilateCNCARMExeName)


C2: MakeBuildDir
	cd C2C_Server && $(assimX64BuildDirections) go build -o $(assimilateCNCx64ExeName)
	cd C2C_Server && mv $(assimilateCNCx64ExeName) ../build/$(assimilateCNCx64ExeName)

fullx64: Victim Prox C2

fullARM: ARMVictim ARMProx ARMC2
