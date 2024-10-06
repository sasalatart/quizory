import { useEffect, useRef } from 'react';
import { Canvas, useFrame, useLoader } from '@react-three/fiber';
import * as THREE from 'three';
import skyboxImage from './skybox.jpg';

interface Props {
  children: React.ReactNode;
}

export function WithSkybox({ children }: Props): JSX.Element {
  // Calculate and set the CSS variable '--vh' to 1% of the current viewport height.
  // This accounts for mobile browsers (e.g. Chrome on Android) where the viewport height changes
  // due to the address bar appearing or disappearing, ensuring the layout adjusts dynamically.
  useEffect(() => {
    const handleResize = () => {
      const vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty('--vh', `${vh}px`);
    };

    handleResize();

    window.addEventListener('resize', handleResize);
    window.addEventListener('orientationchange', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      window.removeEventListener('orientationchange', handleResize);
    };
  }, []);

  return (
    <div className="h-[calc(var(--vh,1vh)*100)]">
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
