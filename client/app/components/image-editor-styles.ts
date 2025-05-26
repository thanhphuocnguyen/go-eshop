// Helper function to create styles for react-easy-crop
export function createStyles() {
  // Only create styles once
  if (!document.getElementById('react-easy-crop-styles')) {
    const style = document.createElement('style');
    style.id = 'react-easy-crop-styles';
    style.innerHTML = `
      .reactEasyCrop_Container {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        overflow: hidden;
        user-select: none;
        touch-action: none;
        cursor: move;
      }
      
      .reactEasyCrop_Image,
      .reactEasyCrop_Video {
        max-width: 100%;
        max-height: 100%;
        margin: auto;
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        will-change: transform;
      }
      
      .reactEasyCrop_CropArea {
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        border: 1px solid rgba(255, 255, 255, 0.5);
        box-sizing: border-box;
        box-shadow: 0 0 0 9999em rgba(0, 0, 0, 0.5);
        overflow: hidden;
        border-radius: 4px;
      }
    `;
    document.head.appendChild(style);
  }
}
