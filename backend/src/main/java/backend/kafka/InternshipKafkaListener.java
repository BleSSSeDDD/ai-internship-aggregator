package backend.kafka;

import com.aggregator.internship.CompanyInternship;
import backend.service.InternshipService;
import com.google.protobuf.InvalidProtocolBufferException;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class InternshipKafkaListener {

    private final InternshipService internshipService;

    @KafkaListener(topics = "internships", groupId = "internship-db-consumer")
    public void listen(byte[] payload) {
        try {
            CompanyInternship internship = CompanyInternship.parseFrom(payload);
            internshipService.saveFromProto(internship);
        } catch (InvalidProtocolBufferException e) {
            throw new RuntimeException("Ошибка парсинга protobuf", e);
        }
    }
}
