# Count-Min Tree Sketch
[Count-Min Tree Sketch: Approximate counting for NLP](https://arxiv.org/pdf/1604.05492.pdf) - 
Guillaume Pitel, 
Geoffroy Fouquier, 
Emmanuel Marchand, 
Abdul Mouhamadsultane, 

## Abstract
The Count-Min Sketch is a widely adopted structure for approximate event counting in large scale processing. In previous works, the original version of the Count-Min-Sketch (CMS) with conservative update has been improved using approximate counters instead of linear counters. These structures are computationaly efficient and improve the average relative error (ARE) of a CMS at constant memory footprint. These improvements are well suited for NLP tasks, in which one is interested by the low-frequency items. However, if Log counters allow to improve ARE, they produce a residual error due to the approximation. In this paper, we propose the Count-Min Tree Sketch variant with pyramidal counters,which are focused toward taking advantage of the Zipfian distribution of text data.
