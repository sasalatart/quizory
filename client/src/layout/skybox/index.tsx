import { useRef } from 'react';
import { Canvas, useFrame, useLoader } from '@react-three/fiber';
import * as THREE from 'three';
import skyboxImage from './skybox.jpg';

interface Props {
  children: React.ReactNode;
}

export function WithSkybox({ children }: Props): JSX.Element {
  return (
    <div className="h-screen">
      <Canvas
        className="absolute top-0 left-0 w-full h-full -z-10"
        camera={{ position: [0, 0, 0.1], fov: 75, near: 0.1 }}
      >
        <Skybox />
      </Canvas>
      <div className="absolute top-0 left-0 w-full h-full flex flex-col">{children}</div>
    </div>
  );
}

function Skybox(): JSX.Element {
  const sphereRef = useRef<THREE.Mesh>(null);
  const texture = useLoader(THREE.TextureLoader, skyboxImage);

  useFrame(() => {
    if (sphereRef.current) {
      sphereRef.current.rotation.y += 0.0001; // Adjust the speed as needed
    }
  });

  return (
    <mesh ref={sphereRef}>
      <sphereGeometry args={[500, 60, 40]} />
      <meshBasicMaterial map={texture} side={THREE.BackSide} />
    </mesh>
  );
}
